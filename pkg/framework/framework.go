package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"

	"github.com/snobb/susanin/pkg/middleware"
)

const (
	mGet = iota
	mPut
	mPost
	mDelete
	mPatch
	mHead
	mOptions
	mSize
)

// Framework is a web framework main data structure
type Framework struct {
	prefix  string
	methods [mSize]*Router
	stack   []middleware.Middleware
}

// New is the Framework constructor
func New() *Framework {
	return &Framework{
		stack: make([]middleware.Middleware, 0),
	}
}

// NewWithPrefix is the Framework constructor
func NewWithPrefix(prefix string) *Framework {
	return &Framework{
		stack:  make([]middleware.Middleware, 0),
		prefix: prefix,
	}
}

// Attach adds middleware to the chain
func (fw *Framework) Attach(middlewares ...middleware.Middleware) *Framework {
	fw.stack = append(fw.stack, middlewares...)
	return fw
}

func (fw *Framework) handler(method int, path string, handler http.HandlerFunc) {
	if fw.methods[method] == nil {
		fw.methods[method] = NewRouter()
	}

	rt := fw.methods[method]
	if err := rt.Handle(fw.prefix+path, handler); err != nil {
		panic(err)
	}
}

func (fw *Framework) dispatch(w http.ResponseWriter, r *http.Request) {
	var method int

	switch r.Method {
	case http.MethodGet:
		method = mGet

	case http.MethodPut:
		method = mPut

	case http.MethodPost:
		method = mPost

	case http.MethodDelete:
		method = mDelete

	case http.MethodPatch:
		method = mPatch

	case http.MethodHead:
		method = mHead

	case http.MethodOptions:
		method = mOptions

	default:
		returnError(w, "Invalid REST method", 404)
		return
	}

	rt := fw.methods[method]
	if rt == nil {
		returnError(w, "Method is not found", 404)
		return
	}

	rt.RouterHandler(w, r)
}

// ServeHTTP is the implementation of the http.Handler interface
// It combines the chain and serves the HTTP requests.
func (fw *Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = http.HandlerFunc(fw.dispatch)
	for i := 0; i < len(fw.stack); i++ {
		h = fw.stack[i](h)
	}

	h.ServeHTTP(w, r)
}

// Get adds handler for GET requests
func (fw *Framework) Get(path string, handler http.HandlerFunc) {
	fw.handler(mGet, path, handler)
}

// Put adds handler for PUT requests
func (fw *Framework) Put(path string, handler http.HandlerFunc) {
	fw.handler(mPut, path, handler)
}

// Post adds handler for POST requests
func (fw *Framework) Post(path string, handler http.HandlerFunc) {
	fw.handler(mPost, path, handler)
}

// Delete adds handler for DELETE requests
func (fw *Framework) Delete(path string, handler http.HandlerFunc) {
	fw.handler(mDelete, path, handler)
}

// Patch adds handler for PATCH requests
func (fw *Framework) Patch(path string, handler http.HandlerFunc) {
	fw.handler(mPatch, path, handler)
}

// Head adds handler for PATCH requests
func (fw *Framework) Head(path string, handler http.HandlerFunc) {
	fw.handler(mHead, path, handler)
}

// Options adds handler for PATCH requests
func (fw *Framework) Options(path string, handler http.HandlerFunc) {
	fw.handler(mOptions, path, handler)
}

// Clear clears all handlers for all methods
func (fw *Framework) Clear() {
	for i := 0; i < mSize; i++ {
		fw.methods[i] = nil
	}
}
