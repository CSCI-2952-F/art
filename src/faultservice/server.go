package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/triplewy/art/src/faultservice/genproto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"net"
	"os"
)

var (
	port        string
	kubeconfig  string

	zLogger *zap.Logger
	sugar   *zap.SugaredLogger
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")

	port = "8080"
	if value := os.Getenv("PORT"); value != "" {
		port = value
	}

	zLogger, _ = zap.NewProduction()
	sugar = zLogger.Sugar()
}

func main() {
	flag.Parse()

	defer zLogger.Sync()

	run(port, kubeconfig)
}

// run starts the gRPC server
func run(port, kubeconfig string) {
	sugar.Infof("starting grpc server at :%s", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	svc := newServer(kubeconfig)

	pb.RegisterFaultServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	go srv.Serve(lis)
	select {}
}

type server struct {
	ic *IstioClient
	kc *k8sClient
}

func newServer(kubeconfig string) *server {
	sugar.Info("Setting up istio client...")
	ic, err := NewIstioClient(kubeconfig)
	if err != nil {
		sugar.Fatalf("Error setting up istio client: %v", err)
	}
	sugar.Infof("Setuping up k8s client...")
	kc, err := NewK8sClient(kubeconfig)
	if err != nil {
		sugar.Fatalf("Error setting up k8s client: %v", err)
	}
	return &server{
		ic: ic,
		kc: kc,
	}
}

func (s *server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.ic.ApplyFaultInjection(req.GetSvc(), req.GetUri(), req.GetPercent())
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.ic.DeleteFaultInjection(req.GetSvc())
}
