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

func extractTraceIDAndMetadata(ctx context.Context) (string, metadata.MD, context.Context) {
	var traceID string

	var md metadata.MD

	if m, ok := metadata.FromIncomingContext(ctx); ok {
		md = m.Copy()
		if values := md.Get("x-trace-id"); len(values) > 0 {
			traceID = values[0]
		} else {
			traceID = uuid.New().String()
			md.Set("x-trace-id", traceID)
		}

		ctx = metadata.NewIncomingContext(ctx, md)
	}

	return traceID, md, ctx
}

func extractClientIP(ctx context.Context) string {
	var clientIP string

	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = addr.IP.String()
		} else {
			clientIP = p.Addr.String()
		}
	}

	return clientIP
}

func getResponseSize(resp interface{}) int {
	if respMsg, ok := resp.(proto.Message); ok {
		return proto.Size(respMsg)
	}

	return 0
}

func logGrpcCall(logger *slog.Logger, traceID, clientIP, method string, req interface{}, duration float64, size int, md metadata.MD, err error) {
	if err != nil {
		logger.Warn("grpc call",
			"trace_id", traceID,
			"metadata", md,
			"method", method,
			"req", MaskSensitiveFields(req),
			"duration", duration,
			"IP-Address", clientIP,
			"size", size,
			"error", err,
		)
	} else {
		logger.Info("grpc call",
			"trace_id", traceID,
			"metadata", md,
			"method", method,
			"req", MaskSensitiveFields(req),
			"duration", duration,
			"IP-Address", clientIP,
			"size", size,
		)
	}
}

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		traceID, md, ctx := extractTraceIDAndMetadata(ctx)

		clientIP := extractClientIP(ctx)

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		size := getResponseSize(resp)

		logGrpcCall(logger, traceID, clientIP, info.FullMethod, req, duration, size, md, err)

		return resp, err
	}
}
