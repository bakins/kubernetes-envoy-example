package util

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var zipkinHeaders = []string{
	"x-ot-span-context",
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
}

// CopyZipkinHeaders copies request headers to the request context.
func CopyZipkinHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]string, 7)
		for _, key := range zipkinHeaders {
			v := r.Header.Get(key)
			if v != "" {
				data[key] = v
			}
		}
		md := metadata.New(data)
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// UnaryServerInterceptor returns a new unary server interceptor for
// copying metadata
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		incoming, ok := metadata.FromIncomingContext(ctx)

		if ok {
			outgoing, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				outgoing = metadata.MD{}
			}
			for _, key := range zipkinHeaders {
				v, ok := incoming[key]
				if ok {
					outgoing[key] = v
				}
			}
			ctx = metadata.NewOutgoingContext(ctx, outgoing)
		}

		resp, err := handler(ctx, req)

		return resp, err
	}
}
