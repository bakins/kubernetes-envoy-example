package frontend

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	// import for pprof
	_ "net/http/pprof"

	"github.com/bakins/kubernetes-envoy-example/api/user"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// OptionsFunc sets options when creating a new server.
type OptionsFunc func(*Server) error

// Server is a wrapper for a simple front end HTTP server
type Server struct {
	address    string
	endpoint   string
	auxAddress string
	server     *http.Server
	auxServer  *http.Server
	user       user.UserServiceClient
}

// New creates a new server
func New(options ...OptionsFunc) (*Server, error) {
	s := &Server{
		address:    ":8080",
		endpoint:   "127.0.0.1:9090",
		auxAddress: ":9999",
	}

	for _, f := range options {
		if err := f(s); err != nil {
			return nil, errors.Wrap(err, "options function failed")
		}
	}

	/*
		mw := grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.NoopTracer{})),
			grpc_prometheus.UnaryClientInterceptor,
		)
	*/
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		s.endpoint,
		grpc.WithInsecure(),
		//grpc.WithUnaryInterceptor(mw),
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not create grpc client")
	}

	s.user = user.NewUserServiceClient(conn)

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
	mux.Handle("/", nethttp.Middleware(opentracing.NoopTracer{}, http.HandlerFunc(s.index)))

	s.server = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	s.auxServer = &http.Server{
		Addr:    s.auxAddress,
		Handler: http.DefaultServeMux,
	}

	var g errgroup.Group

	g.Go(func() error {
		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				return errors.Wrap(err, "failed to run main http server")
			}
		}
		return nil
	})

	g.Go(func() error {
		if err := s.auxServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				return errors.Wrap(err, "failed to run aux http server")
			}
		}
		return nil
	})

	return g.Wait()
}

// Stop will stop the server
func (s *Server) Stop() {

	for _, srv := range []*http.Server{s.auxServer, s.server} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		srv.Shutdown(ctx)
	}
}

func healthz(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "OK\n")
}

func (s *Server) index(wr http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

	users, err := s.user.ListUsers(context.Background(), &user.ListUsersRequest{})
	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}
	fmt.Fprintln(wr, users)
}
