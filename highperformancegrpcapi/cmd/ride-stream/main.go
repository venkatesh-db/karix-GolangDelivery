package main

import (
	"context"
	"net"
	"net/http"
	"time"

	ridev1 "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"github.com/example/highperformancegrpcapi/internal/config"
	"github.com/example/highperformancegrpcapi/internal/matching"
	"github.com/example/highperformancegrpcapi/internal/server"
	"github.com/example/highperformancegrpcapi/internal/stream"
	"github.com/example/highperformancegrpcapi/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	_, _ = maxprocs.Set()
	log, _ := zap.NewProduction()
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config load failed", zap.Error(err))
	}

	reg := prometheus.NewRegistry()
	metrics := telemetry.New(reg)

	broker := stream.NewBroker(cfg.ShardCount, cfg.MaxSessions, metrics)
	engine := matching.New(broker, log, 256)
	engine.Start(context.Background())

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             15 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	h := &server.RideStreamHandler{Broker: broker, Engine: engine, Metrics: metrics, Log: log, BufSize: cfg.OutboundBufferSize}
	ridev1.RegisterRideStreamServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	go func() {
		l, err := net.Listen("tcp", cfg.GRPCListenAddr)
		if err != nil {
			log.Fatal("grpc listen failed", zap.Error(err))
		}
		log.Info("grpc listening", zap.String("addr", cfg.GRPCListenAddr))
		if err := grpcServer.Serve(l); err != nil {
			log.Fatal("grpc server exited", zap.Error(err))
		}
	}()

	httpServer := &http.Server{Addr: cfg.MetricsListenAddr, Handler: telemetry.Handler(reg)}
	go func() {
		log.Info("metrics listening", zap.String("addr", cfg.MetricsListenAddr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("metrics server exited", zap.Error(err))
		}
	}()

	select {}
}
