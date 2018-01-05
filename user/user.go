package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bakins/kubernetes-envoy-example/api/user"
	"github.com/bakins/kubernetes-envoy-example/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	//grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func init() {
	grpc_prometheus.EnableHandlingTimeHistogram()
}

// OptionsFunc sets options when creating a new server.
type OptionsFunc func(*Server) error

// Server is a wrapper for a simple front end HTTP server
type Server struct {
	address string
	store   *userStore
	s       *util.Server
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

	s.store = newUserStore(s)
	s.store.LoadSampleData()

	return s, nil
}

// SetAddress sets the listening address.
func SetAddress(address string) OptionsFunc {
	return func(s *Server) error {
		s.address = address
		return nil
	}
}

// Run starts the server. This generally does not return.
func (s *Server) Run() error {
	logger, err := util.NewDefaultLogger()
	if err != nil {
		return errors.Wrapf(err, "failed to create logger")
	}

	grpc_zap.ReplaceGrpcLogger(logger)
	grpc_prometheus.EnableHandlingTimeHistogram()

	//tracer := util.NewTracer("user")

	g := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				util.UnaryServerInterceptor(),
				util.UnaryServerSleeperInterceptor(time.Second*3),
				grpc_validator.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)

	user.RegisterUserServiceServer(g, s.store)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", healthz)

	h := &http.Server{
		Handler: mux,
	}

	s.s = util.NewServer(g, s.address, h, ":11111")

	if err := s.s.Run(); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "failed to start server")
		}
	}

	return nil
}

// Stop will stop the server
func (s *Server) Stop() {
	s.s.Stop()
}

func healthz(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "OK\n")
}
