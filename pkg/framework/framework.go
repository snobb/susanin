package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
	"path"

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

// Route callback function
type Route func()

// Framework is a web framework main data structure
type Framework struct {
	methods         [mSize]*Router
	middlewares     []middleware.Middleware
	prefixes        []string
	notFoundHandler http.Handler
}

// New is the Framework constructor
func New() *Framework {
	return &Framework{}
}

// WithDefaultPrefix adds default prefix for all following route declarations. All previous ones
// will NOT have the default prefix assigned.
// As an alternative http.StripPrefix can be used to wrap the main Framework instance. That will
// have the same effect as setting the prefix with this method right after creating the instance.
func (fw *Framework) WithDefaultPrefix(prefix string) *Framework {
	fw.prefixes = append(fw.prefixes, prefix)
	return fw
}

// WithPrefix registers paths with given prefix.
func (fw *Framework) WithPrefix(prefix string, route Route) *Framework {
	fw.prefixes = append(fw.prefixes, prefix)
	defer func() {
		fw.prefixes = fw.prefixes[:len(fw.prefixes)-1]
	}()

	route()
	return fw
}

// WithNotFoundHandler sets NotFoundHander that will be used in case the route is not found.
func (fw *Framework) WithNotFoundHandler(notFoundHander http.Handler) *Framework {
	fw.notFoundHandler = notFoundHander
	return fw
}

// Attach adds middleware to the chain
func (fw *Framework) Attach(middlewares ...middleware.Middleware) *Framework {
	fw.middlewares = append(fw.middlewares, middlewares...)
	return fw
}

func (fw *Framework) handler(method int, pattern string, handler http.Handler) {
	if fw.methods[method] == nil {
		fw.methods[method] = NewRouter(fw.notFoundHandler)
	}

	rt := fw.methods[method]

	pp := append([]string{}, fw.prefixes...)
	pp = append(pp, pattern)

	if err := rt.Handle(path.Join(pp...), handler); err != nil {
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
	for i := 0; i < len(fw.middlewares); i++ {
		h = fw.middlewares[i](h)
	}

	h.ServeHTTP(w, r)
}

// Get adds handler for GET requests
func (fw *Framework) Get(path string, handler http.Handler) {
	fw.handler(mGet, path, handler)
}

// Put adds handler for PUT requests
func (fw *Framework) Put(path string, handler http.Handler) {
	fw.handler(mPut, path, handler)
}

// Post adds handler for POST requests
func (fw *Framework) Post(path string, handler http.Handler) {
	fw.handler(mPost, path, handler)
}

// Delete adds handler for DELETE requests
func (fw *Framework) Delete(path string, handler http.Handler) {
	fw.handler(mDelete, path, handler)
}

// Patch adds handler for PATCH requests
func (fw *Framework) Patch(path string, handler http.Handler) {
	fw.handler(mPatch, path, handler)
}

// Head adds handler for PATCH requests
func (fw *Framework) Head(path string, handler http.Handler) {
	fw.handler(mHead, path, handler)
}

// Options adds handler for PATCH requests
func (fw *Framework) Options(path string, handler http.Handler) {
	fw.handler(mOptions, path, handler)
}

// Clear clears all handlers for all methods
func (fw *Framework) Clear() {
	for i := 0; i < mSize; i++ {
		fw.methods[i] = nil
	}
}
