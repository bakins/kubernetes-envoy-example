package user

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// OptionsFunc sets options when creating a new server.
type OptionsFunc func(*Server) error

// Server is a wrapper for a simple front end HTTP server
type Server struct {
	address string
	server  *http.Server
}

// New creates a new server
func New(options ...OptionsFunc) (*Server, error) {
	s := &Server{
		address: ":8080",
	}

	for _, f := range options {
		if err := f(s); err != nil {
			return nil, errors.Wrap(err, "options function failed")
		}
	}

	return s, nil
}

// SetAddress sets the listening address.
func SetAddress(address string) OptionsFunc {
	return func(s *Server) error {
		s.address = address
		return nil
	}
}

// Run starts the HTTP server. THis generally does not return.
func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/", index)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: nethttp.Middleware(opentracing.NoopTracer{}, mux),
	}

	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "failed to run http server")
		}
	}

	return nil
}

// Stop will stop the server
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.server.Shutdown(ctx)
}

func healthz(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "OK\n")
}

func index(wr http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
	fmt.Fprintf(wr, "Hello\n")
}
