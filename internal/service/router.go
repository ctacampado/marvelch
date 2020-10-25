package service

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Mux is an abstraction for a router/mux implementation
// so that we can easily replace our router with
// any 3rd party routers or even our own implementation
type Mux interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// SvcRouter is the service router. The Init field is a function
// that sets the different routes/endpoints
type SvcRouter struct {
	Mux      Mux
	initFunc func(Mux)
}

// InitRoutes is a wrapper method for the InitFunc
func (r *SvcRouter) InitRoutes() {
	r.Mux = chi.NewRouter()
	r.initFunc(r.Mux)
}
