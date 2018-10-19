package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
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

// Middleware is a type for Middleware function
type Middleware func(http.Handler) http.Handler

// Framework is a web framework main data structure
type Framework struct {
	methods [mSize]*Router
	stack   []Middleware
}

// NewFramework is the Framework constructor
func NewFramework() *Framework {
	return &Framework{
		stack: make([]Middleware, 0),
	}
}

// Attach adds middleware to the chain
func (s *Framework) Attach(middlewares ...Middleware) *Framework {
	s.stack = append(s.stack, middlewares...)
	return s
}

func error404(w http.ResponseWriter, msg string) {
	http.Error(w, msg, 404)
}

func (s *Framework) handler(method int, path string, handler http.HandlerFunc) error {
	if s.methods[method] == nil {
		s.methods[method] = NewRouter()
	}

	rt := s.methods[method]
	return rt.Handle(path, handler)
}

func (s *Framework) dispatch(w http.ResponseWriter, r *http.Request) {
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
		error404(w, "Invalid REST method")
		return
	}

	rt := s.methods[method]
	if rt == nil {
		error404(w, "Not found")
		return
	}

	rt.RouterHandler(w, r)
}

// Router combines the chain and returns the resulting handler function
func (s *Framework) Router() http.Handler {
	var h http.Handler = http.HandlerFunc(s.dispatch)
	for i := 0; i < len(s.stack); i++ {
		h = s.stack[i](h)
	}

	return h
}

// Get adds handler for GET requests
func (s *Framework) Get(path string, handler http.HandlerFunc) {
	s.handler(mGet, path, handler)
}

// Put adds handler for PUT requests
func (s *Framework) Put(path string, handler http.HandlerFunc) {
	s.handler(mPut, path, handler)
}

// Post adds handler for POST requests
func (s *Framework) Post(path string, handler http.HandlerFunc) {
	s.handler(mPost, path, handler)
}

// Delete adds handler for DELETE requests
func (s *Framework) Delete(path string, handler http.HandlerFunc) {
	s.handler(mDelete, path, handler)
}

// Patch adds handler for PATCH requests
func (s *Framework) Patch(path string, handler http.HandlerFunc) {
	s.handler(mPatch, path, handler)
}

// Head adds handler for PATCH requests
func (s *Framework) Head(path string, handler http.HandlerFunc) {
	s.handler(mHead, path, handler)
}

// Options adds handler for PATCH requests
func (s *Framework) Options(path string, handler http.HandlerFunc) {
	s.handler(mOptions, path, handler)
}

// Clear clears all handlers for all methods
func (s *Framework) Clear() {
	for i := 0; i < mSize; i++ {
		s.methods[i] = nil
	}
}
