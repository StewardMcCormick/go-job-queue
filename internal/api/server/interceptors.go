package server

import (
	"context"
	"fmt"
	"time"

	appctx "github.com/StewardMcCormick/go-job-queue/internal/api/app_context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnariRequestIdInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		var requestId string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md["x-request-id"]) > 0 {
			requestId = md["x-request-id"][0]
		} else {
			requestId = uuid.New().String()
		}

		log := log.With(zap.String("request_id", requestId))
		ctx = appctx.WithLogger(ctx, log)

		if p, ok := peer.FromContext(ctx); ok {
			log.Info(fmt.Sprintf("[NEW REQUEST] From %s to %s", p.Addr.String(), info.FullMethod))
		}

		resp, err := handler(ctx, req)
		if err != nil {
			log.Info(fmt.Sprintf("[REQUEST COMPLETED WITH ERROR] %s, total duration: %d ms",
				err, time.Since(start).Milliseconds()))
		} else {
			log.Info(fmt.Sprintf("[REQUEST COMPLETED] total duration: %d ms", time.Since(start).Milliseconds()))
		}

		return resp, err
	}
}
