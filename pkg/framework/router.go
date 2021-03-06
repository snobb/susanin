package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type valuesKey struct{}

const rootLink = "#ROOT#"

// Router is a URI path router object
type Router struct {
	root *chainLink
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
func NewRouter() *Router {
	return &Router{
		root: newChainLink(rootLink),
	}
}

// Handle add a route and a handler
func (rt *Router) Handle(path string, handler http.Handler) (err error) {
	splatIdx := strings.IndexRune(path, '*')

	if splatIdx != -1 && splatIdx != len(path)-1 {
		return fmt.Errorf("invalid path: splat must be at the end of the path")
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
				return fmt.Errorf("conflict: duplicate pattern at the same level")
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
		return fmt.Errorf("handler already exists")
	}

	cur.handler = handler

	return
}

// Lookup for a handler in the path, a handler, pattern values and error is returned.
func (rt *Router) Lookup(path string) (http.Handler, map[string]string, error) {
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
	hasSplat := false
	found := false

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
		return cur.handler, values, nil
	}

	if hasSplat {
		return splatHandler, values, nil
	}

	return nil, nil, fmt.Errorf("Endpoint is not found")
}

// RouterHandler is a http.HandlerFunc router that dispatches the request
// based on saved routes and handlers
func (rt *Router) RouterHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	handler, values, err := rt.Lookup(uri)
	if err != nil {
		returnError(w, err.Error(), 404)
		return
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
	body, _ := json.Marshal(map[string]interface{}{
		"code": code,
		"msg":  msg,
	})

	w.WriteHeader(code)
	w.Write(body)
}
