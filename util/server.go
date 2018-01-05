package util

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Server is a wrapper for a grpc and an http server
type Server struct {
	httpAddress string
	httpServer  *http.Server
	grpcAddress string
	grpcServer  *grpc.Server
}

func NewServer(g *grpc.Server, gAddr string, h *http.Server, hAddr string) *Server {
	return &Server{
		httpAddress: hAddr,
		httpServer:  h,
		grpcAddress: gAddr,
		grpcServer:  g,
	}
}

// Stop ...
func (s *Server) Stop() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.grpcServer.GracefulStop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_ = s.httpServer.Shutdown(ctx)
	}()

	wg.Wait()
}

// Run generally does not return
func (s *Server) Run() error {
	errChan := make(chan error)

	go func() {
		errChan <- s.startGRPCServer()
	}()

	go func() {
		errChan <- s.startHTTPServer()
	}()

	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			// if the first message is an error, we need to
			// receive on the channel to avoid leaking the channel and/or error.
			// this probably doesn't matter as the server will usually exit
			// if we return an error
			if i == 0 {
				go func() {
					<-errChan
				}()
			}
			return err
		}
	}

	return nil
}

func (s *Server) startGRPCServer() error {
	l, err := net.Listen("tcp", s.grpcAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", s.grpcAddress)
	}

	if err := s.grpcServer.Serve(l); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrapf(err, "failed to start grpc server", s.grpcAddress)
		}
	}
	return nil
}

func (s *Server) startHTTPServer() error {
	l, err := net.Listen("tcp", s.httpAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", s.httpServer)
	}

	if err := s.httpServer.Serve(l); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "failed to start http server")
		}
	}
	return nil
}
