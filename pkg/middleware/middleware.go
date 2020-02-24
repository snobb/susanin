package middleware

import (
	"net/http"
)

// Middleware is a type for Middleware function
type Middleware func(http.Handler) http.Handler
