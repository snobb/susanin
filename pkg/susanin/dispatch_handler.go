package susanin

/**
 * @author: Alex Kozadaev
 */

import "net/http"

// MiddleWare is a type for MiddleWare function
type MiddleWare func(http.HandlerFunc) http.HandlerFunc

// DispatchHandler is a struct that chain middleware
type DispatchHandler struct {
	stack   []MiddleWare
	handler http.HandlerFunc
}

// Attach adds middleware to the chain
func (s *DispatchHandler) Attach(next MiddleWare) *DispatchHandler {
	s.stack = append(s.stack, next)
	return s
}

// Handler add the handler to the chain and get the resulting handler function
func (s *DispatchHandler) Handler(h http.HandlerFunc) http.HandlerFunc {
	for i := 0; i < len(s.stack); i++ {
		h = s.stack[i](h)
	}
	s.handler = h
	return s.handler
}
