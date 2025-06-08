package middleware

import (
	"net/http"
)

// Middleware представляет собой функцию-обертку для http.Handler
type Middleware func(http.Handler) http.Handler

// Chain объединяет несколько middleware в одну цепочку
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
