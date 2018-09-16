package susanin

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type valueKey string

const valuesKey valueKey = "values"
const rootLink = "#ROOT#"

// Susanin is a URI path router object
type Susanin struct {
	root *chainLink
}

type chainLink struct {
	name      string
	nextConst map[string]*chainLink
	nextVar   *chainLink
	nextSplat *chainLink
	handler   http.HandlerFunc
}

func newChainLink(token string) *chainLink {
	// strip colon from variable name
	if token[0] == ':' {
		token = token[1:]
	}

	return &chainLink{
		name: token,
	}
}

// NewSusanin creates a new Susanin instance
func NewSusanin() *Susanin {
	return &Susanin{
		root: newChainLink(rootLink),
	}
}

// Handle add a route and a handler
func (s *Susanin) Handle(path string, handler http.HandlerFunc) (err error) {
	splatIdx := strings.IndexRune(path, '*')

	if splatIdx != -1 && splatIdx != len(path)-1 {
		return errors.New("invalid path: splat must be at the end of the path")
	}

	if path[0] == '/' {
		path = path[1:]
	}

	// clear trailing slash so that it matches /path and /path/
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	tokens := strings.Split(path, "/")

	cur := s.root

	for _, token := range tokens {
		switch {
		case token[0] == ':': // variable
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

	return
}

// Lookup for a handler in the path, a handler, pattern values and error is returned.
func (s *Susanin) Lookup(path string) (http.HandlerFunc, map[string]string, error) {
	if path[0] == '/' {
		path = path[1:]
	}

	// clear trailing slash so that it matches /path and /path/
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	tokens := strings.Split(path, "/")

	cur := s.root
	var splatHandler http.HandlerFunc
	var values map[string]string
	hasSplat := false
	found := false

	for _, token := range tokens {
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

		if cur.nextSplat != nil {
			splatHandler = cur.nextSplat.handler
			hasSplat = true
		}

		if !found {
			break
		}
	}

	if found {
		return cur.handler, values, nil
	}

	if hasSplat {
		return splatHandler, values, nil
	}

	return nil, nil, errors.New("not found")
}

// Router is a http.HandlerFunc router that dispatches the request
// based on saved routes and handlers
func (s *Susanin) Router(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	handler, values, err := s.Lookup(uri)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	if len(values) > 0 {
		ctx := r.Context()
		ctx = context.WithValue(ctx, valuesKey, values)
		r = r.WithContext(ctx)
	}

	handler(w, r)
}

// GetValues gets the values from match patters from the http.Request context
func GetValues(r *http.Request) (map[string]string, bool) {
	ctx := r.Context()
	value := ctx.Value(valuesKey)
	if value == nil {
		return nil, false
	}

	values, ok := value.(map[string]string)
	return values, ok
}
