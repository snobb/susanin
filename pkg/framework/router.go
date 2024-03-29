package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type valuesKey struct{}

const rootLink = "#ROOT#"

// Router is a URI path router object
type Router struct {
	root            *chainLink
	notFoundHandler http.Handler
}

type chainLink struct {
	name      string
	nextConst map[string]*chainLink
	nextVar   *chainLink
	nextSplat *chainLink
	handler   http.Handler
}

func newChainLink(token string) *chainLink {
	// strip colon from variable name
	if len(token) > 0 && token[0] == ':' {
		token = token[1:]
	}

	return &chainLink{
		name: token,
	}
}

// NewRouter creates a new Router instance
func NewRouter(notFoundHandler http.Handler) *Router {
	return &Router{
		root:            newChainLink(rootLink),
		notFoundHandler: notFoundHandler,
	}
}

// Handle add a route and a handler
func (rt *Router) Handle(path string, handler http.Handler) (err error) {
	splatIdx := strings.IndexRune(path, '*')

	if splatIdx != -1 && splatIdx != len(path)-1 {
		return errors.New("invalid path: splat must be at the end of the path")
	}

	if path[0] == '/' {
		path = path[1:]
	}

	// clear trailing slash so that it matches /path and /path/
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	tokens := strings.Split(path, "/")

	cur := rt.root

	for _, token := range tokens {
		switch {
		case len(token) > 0 && token[0] == ':': // variable
			if cur.nextVar == nil {
				cur.nextVar = newChainLink(token)
			} else if token[1:] != cur.nextVar.name {
				return errors.New("conflict: duplicate pattern at the same level")
			}

			cur = cur.nextVar

		case token == "*": // splat
			cur.nextSplat = newChainLink(token)
			cur = cur.nextSplat

		default:
			var next *chainLink
			if cur.nextConst == nil {
				cur.nextConst = make(map[string]*chainLink)
				next = newChainLink(token)
				cur.nextConst[token] = next
			} else {
				if found, ok := cur.nextConst[token]; ok {
					next = found
				} else {
					next = newChainLink(token)
					cur.nextConst[token] = next
				}
			}

			cur = next
		}
	}

	if cur.handler != nil {
		return errors.New("handler already exists")
	}

	cur.handler = handler

	return nil
}

// Lookup for a handler in the path, a handler and pattern values is returned.
// If handler is not found the function returns NotFoundHandler configured for the router (can be
// nil).
func (rt *Router) Lookup(path string) (http.Handler, map[string]string) {
	if path[0] == '/' {
		path = path[1:]
	}

	// clear trailing slash so that it matches /path and /path/
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	tokens := strings.Split(path, "/")

	cur := rt.root
	var splatHandler http.Handler
	var values map[string]string
	hasSplat, found := false, false

	for _, token := range tokens {
		if cur.nextSplat != nil {
			splatHandler = cur.nextSplat.handler
			hasSplat = true
		}

		found = false
		if next, ok := cur.nextConst[token]; ok {
			cur = next
			found = true
		} else if cur.nextVar != nil {
			cur = cur.nextVar
			if values == nil {
				values = make(map[string]string)
			}
			values[cur.name] = token
			found = true
		}

		if !found {
			break
		}
	}

	if found && cur.handler != nil {
		return cur.handler, values
	}

	if hasSplat {
		return splatHandler, values
	}

	return rt.notFoundHandler, nil
}

// RouterHandler is a http.HandlerFunc router that dispatches the request
// based on saved routes and handlers
func (rt *Router) RouterHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	handler, values := rt.Lookup(uri)
	if handler == nil {
		// set default NotFoundHandler
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			returnError(w, "Endpoint is not found", 404)
		})
	}

	if len(values) > 0 {
		ctx := r.Context()
		ctx = context.WithValue(ctx, valuesKey{}, values)
		r = r.WithContext(ctx)
	}

	handler.ServeHTTP(w, r)
}

// GetValues gets the match pattern values from the http.Request context
func GetValues(ctx context.Context) (map[string]string, bool) {
	value := ctx.Value(valuesKey{})
	if value == nil {
		return nil, false
	}

	values, ok := value.(map[string]string)
	return values, ok
}

func returnError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{code, msg})
}
