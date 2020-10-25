package service

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Service contains common elements
// for a microservice
type Service struct {
	SvcName   string
	SvcPort   string
	SvcHost   string
	SvcRouter SvcRouter
	SvcCache  SvcCache
}

// Start starts the service
func (s *Service) Start(msg string) error {
	s.SvcName = os.Getenv("SVCNAME")
	s.SvcHost = os.Getenv("HOST")
	s.SvcPort = os.Getenv("PORT")
	s.SvcRouter.InitRoutes()
	log.Println(s.SvcName + " " + msg + s.SvcHost + ":" + s.SvcPort)
	log.Printf("Service Start svc %+v\n", s)
	return http.ListenAndServe(s.SvcHost+":"+s.SvcPort, s.SvcRouter.Mux)
}

// Builder is for building the
// Service struct with builder pattern
type Builder struct {
	s Service
}

// Build returns the built service
func (b *Builder) Build() Service {
	return b.s
}

// LoadEnv loads .env to ENV Vars
func (b *Builder) LoadEnv() *Builder {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return b
}

// Cache sets the service cache
func (b *Builder) Cache(f func(Cache) error) *Builder {
	b.s.SvcCache.InitCache(f)
	return b
}

// Router sets initialization function for the router
func (b *Builder) Router(f func(Mux)) *Builder {
	b.s.SvcRouter.initFunc = f
	return b
}
