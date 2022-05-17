// Package middleware provides middleware methods.
package middleware

import (
	"net/http"
)

// Middleware defines a type for usage in middleware conveyor.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Conveyor acts as a conveyor of one ttp.HandlerFunc to multiple middlewares.
func Conveyor(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
