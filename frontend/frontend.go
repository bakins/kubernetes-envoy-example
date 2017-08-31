package frontend

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/bakins/kubernetes-envoy-example/api/order"
	"github.com/bakins/kubernetes-envoy-example/api/user"
	"github.com/bakins/kubernetes-envoy-example/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	basictracer "github.com/opentracing/basictracer-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// OptionsFunc sets options when creating a new server.
type OptionsFunc func(*Server) error

// Server is a wrapper for a simple front end HTTP server
type Server struct {
	address  string
	endpoint string
	server   *http.Server
	user     user.UserServiceClient
	order    order.OrderServiceClient
}

type noopRecorder struct{}

func (n noopRecorder) RecordSpan(basictracer.RawSpan) {}

// New creates a new server
func New(options ...OptionsFunc) (*Server, error) {
	s := &Server{
		address:  ":8080",
		endpoint: "127.0.0.1:9090",
	}

	for _, f := range options {
		if err := f(s); err != nil {
			return nil, errors.Wrap(err, "options function failed")
		}
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		s.endpoint,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_prometheus.UnaryClientInterceptor,
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_prometheus.StreamClientInterceptor,
		)),
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not create grpc client")
	}

	// use the same connection for each. Envoy will handle
	// load balancing, etc
	s.user = user.NewUserServiceClient(conn)
	s.order = order.NewOrderServiceClient(conn)
	return s, nil
}

// SetAddress sets the listening address.
func SetAddress(address string) OptionsFunc {
	return func(s *Server) error {
		s.address = address
		return nil
	}
}

// SetEndpoint sets the address for contacting other services.
func SetEndpoint(address string) OptionsFunc {
	return func(s *Server) error {
		s.endpoint = address
		return nil
	}
}

// Run starts the HTTP server. This generally does not return.
func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	//mux.Handle("/index", nethttp.Middleware(s.tracer, http.HandlerFunc(s.index)))
	mux.Handle("/index", util.CopyZipkinHeaders(http.HandlerFunc(s.index)))

	s.server = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "failed to run main http server")
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

func (s *Server) index(wr http.ResponseWriter, r *http.Request) {
	users, err := s.user.ListUsers(r.Context(), &user.ListUsersRequest{})

	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}

	time.Sleep(time.Duration(rand.Intn(3000)) * time.Microsecond)

	orders, err := s.order.ListOrders(r.Context(), &order.ListOrdersRequest{})

	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}

	fmt.Fprintln(wr, users)
	fmt.Fprintln(wr, orders)
}
