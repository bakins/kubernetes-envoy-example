package util

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// UnaryServerSleeperInterceptor returns a new unary server interceptor
// that sleeps for a random amount of time
func UnaryServerSleeperInterceptor(t time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//time.Sleep(time.Duration(rand.Intn(int(t.Nanoseconds()))) * time.Nanosecond)
		return handler(ctx, req)
	}
}
