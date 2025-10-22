package main

import (
	"cards-service/internal/adapters/api"
	"cards-service/internal/config"
	"cards-service/internal/core/app"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	reqvalidator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mwinyimoha/commons/pkg/logging"
	"github.com/mwinyimoha/protos/gen/go/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger, err := logging.NewLoggerConfig().BuildLogger()
	if err != nil {
		log.Fatal("could not initialize logging", err)
	}

	defer logger.Sync()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("could not initialize app config", zap.String("original_error", err.Error()))
	}

	svc, err := app.NewService()
	if err != nil {
		logger.Fatal("could not initialize service", zap.String("original_error", err.Error()))
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.Fatal("could not initialize request validator", zap.String("original_error", err.Error()))
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclogging.UnaryServerInterceptor(api.RequestLogInterceptor(logger)),
			reqvalidator.UnaryServerInterceptor(validator),
			recovery.UnaryServerInterceptor(),
		),
	)

	srv := api.NewServer(svc)
	pb.RegisterCardsServiceServer(s, srv)

	healthsrv := health.NewServer()
	healthpb.RegisterHealthServer(s, healthsrv)

	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.ServerPort))
	if err != nil {
		logger.Fatal("could not bind port", zap.String("original_error", err.Error()))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func(c chan os.Signal) {
		logger.Info(
			"starting gRPC server",
			zap.String("service_name", cfg.ServiceName),
			zap.String("service_version", cfg.ServiceVersion),
			zap.Int("port", cfg.ServerPort),
		)

		if err := s.Serve(lis); err != nil {
			logger.Error("could not start server", zap.String("original_error", err.Error()))
			c <- syscall.SIGTERM
		}
	}(ch)

	received := <-ch

	func() {
		logger.Info("initiating graceful shutdown", zap.String("os_signal", received.String()))

		s.GracefulStop()

	}()
}
