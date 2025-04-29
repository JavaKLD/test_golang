package middleware

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		traceID := ""
		clientIP := ""
		var md metadata.MD

		if m, ok := metadata.FromIncomingContext(ctx); ok {
			md = m.Copy()
			slog.Info("gRPC metadata received", "metadata", md)

			if values := md.Get("x-trace-id"); len(values) > 0 {
				traceID = values[0]
			} else {
				traceID = uuid.New().String()
				md.Set("x-trace-id", traceID)
			}

			ctx = metadata.NewIncomingContext(ctx, md)
		}

		if p, ok := peer.FromContext(ctx); ok {
			if addr, ok := p.Addr.(*net.TCPAddr); ok {
				clientIP = addr.IP.String()
			} else {
				clientIP = p.Addr.String()
			}
		}

		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		if err != nil {
			logger.Warn("grpc call",
				"trace_id", traceID,
				"metadate", md,
				"method", info.FullMethod,
				"req", MaskSensitiveFields(req),
				"duration", duration,
				"IP-Adress", clientIP,
				"size", proto.Size(resp.(proto.Message)),
				"error", err,
			)
		} else {
			logger.Info("grpc call",
				"trace_id", traceID,
				"metadate", md,
				"method", info.FullMethod,
				"req", MaskSensitiveFields(req),
				"duration", duration,
				"IP-Adress", clientIP,
				"size", proto.Size(resp.(proto.Message)),
			)
		}

		return resp, err
	}
}
