package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/bakins/kubernetes-envoy-example/api/item"
	"github.com/bakins/kubernetes-envoy-example/api/order"
	"github.com/bakins/kubernetes-envoy-example/api/user"
	"github.com/bakins/kubernetes-envoy-example/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
	address  string
	endpoint string
	server   *http.Server
	user     user.UserServiceClient
	order    order.OrderServiceClient
	item     item.ItemServiceClient
	logger   *zap.Logger
}

// New creates a new server
func New(options ...OptionsFunc) (*Server, error) {
	s := &Server{
		address:  ":8080",
		endpoint: "127.0.0.1:9090",
	}

	l, err := util.NewDefaultLogger()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create logger")
	}

	s.logger = l

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
	s.item = item.NewItemServiceClient(conn)
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
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/index", newLogMiddleware(util.CopyZipkinHeaders(http.HandlerFunc(s.index)), s.logger))

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

type orderReport struct {
	User   *user.User                      `json:"user"`
	Orders []*order.GetOrderDetailResponse `json:"orders"`
}

func (s *Server) index(wr http.ResponseWriter, r *http.Request) {

	// TODO: we should use a deadline context
	ctx := r.Context()

	report := []*orderReport{}

	users, err := s.user.ListUsers(ctx, &user.ListUsersRequest{})

	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}

	// get orders for each user
	for _, u := range users.Users {
		orders, err := s.order.ListOrders(ctx, &order.ListOrdersRequest{User: u.Id})

		rep := &orderReport{
			User:   u,
			Orders: []*order.GetOrderDetailResponse{},
		}

		if err != nil {
			fmt.Println("failed to list orders for", u.Id, err)
			http.Error(wr, err.Error(), 500)
			return
		}

		// get details
		for _, i := range orders.Orders {
			details, err := s.order.GetOrderDetail(ctx, &order.GetOrderDetailRequest{Id: i.Id})
			if err != nil {
				fmt.Println("failed to get order details for", u.Id, err)
				http.Error(wr, err.Error(), 500)
				return
			}
			rep.Orders = append(rep.Orders, details)
		}
		report = append(report, rep)
	}

	data, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		fmt.Println("failed to mrshal json", err)
		http.Error(wr, err.Error(), 500)
		return
	}

	wr.Write(data)
}
