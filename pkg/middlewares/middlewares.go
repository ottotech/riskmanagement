package middlewares

import (
	"github.com/ottotech/riskmanagement/pkg/listing"
	"net/http"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc

// MediaPathRequired ensures that there is a valid path on the user's
// pc where all risk matrix pictures can be stored
func MediaPathRequired(lister listing.Service) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := lister.GetMediaPath()
			if err != nil {
				http.Redirect(w, r, "/set-media-path", http.StatusSeeOther)
				return
			}
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
