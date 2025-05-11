package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	duration := time.Since(startTime)
	// log the request
	log.Info().Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Msg("received a request")

	result, err := handler(ctx, req)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	// log the response
	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Str("duration", duration.String()).
		Msg("finished handling request")

	return result, err
}

type ResponeRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponeRecorder) Write(b []byte) (int, error) {
	r.Body = b
	return r.ResponseWriter.Write(b)
}

func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		recorder := &ResponeRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(recorder, req)
		duration := time.Since(start)

		logger := log.Info()
		if recorder.StatusCode != http.StatusOK {
			logger = log.Error().Str("body", string(recorder.Body))
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Int("status_code", recorder.StatusCode).
			Str("status_text", http.StatusText(recorder.StatusCode)).
			Str("duration", duration.String()).
			Msg("finished handling request")
	})
}
