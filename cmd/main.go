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
	"github.com/go-playground/validator/v10"
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
		log.Fatal("could not initialize logging:", err)
	}

	defer func() { _ = logger.Sync() }()

	val := validator.New()
	cfg, err := config.New(val)
	if err != nil {
		logger.Fatal("could not initialize configuration", zap.Error(err))
	}

	svc, err := app.NewService(val)
	if err != nil {
		logger.Fatal("could not initialize service", zap.Error(err))
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.Fatal("could not initialize request validator", zap.Error(err))
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

	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(s, healthSrv)

	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.ServerPort))
	if err != nil {
		logger.Fatal("could not bind port", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info(
			"starting gRPC server",
			zap.String("service_name", cfg.ServiceName),
			zap.String("service_version", cfg.ServiceVersion),
			zap.Int("server_port", cfg.ServerPort),
			zap.Bool("debug_mode", cfg.Debug),
		)

		healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		if err := s.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case signal := <-sigCh:
		logger.Info("initiating graceful shutdown", zap.String("signal", signal.String()))

		healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		s.GracefulStop()
	case err = <-errCh:
		logger.Error("server stopped unexpectedly", zap.Error(err))

		healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		s.Stop()
	}

	logger.Info("server stopped")
}
