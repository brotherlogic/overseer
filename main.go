package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pb "github.com/brotherlogic/overseer/proto"
	pspb "github.com/brotherlogic/pstore/proto"

	ghbclient "github.com/brotherlogic/githubridge/client"
	pstoreclient "github.com/brotherlogic/pstore/client"
)

const (
	CONFIG_KEY = "github.com/brotherlogic/overseer/config"
)

var (
	grpc_port    = flag.Int("grpc_port", 8080, "gRPC port")
	metrics_port = flag.Int("metrics_port", 8082, "Metrics port")
)

type Server struct {
	Client  ghbclient.GithubridgeClient
	Pclient pstoreclient.PStoreClient
}

func (s *Server) loadConfig(ctx context.Context) (*pb.Config, error) {
	// Load the config
	val, err := s.Pclient.Read(ctx, &pspb.ReadRequest{
		Key: CONFIG_KEY,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("Failure to read config: %w", err)
	}

	config := &pb.Config{}
	err = proto.Unmarshal(val.GetValue().GetValue(), config)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal config: %w", err)
	}
	return config, err
}

func (s Server) saveConfig(ctx context.Context, config *pb.Config) error {
	data, err := proto.Marshal(config)
	if err != nil {
		return err
	}

	_, err = s.Pclient.Write(ctx, &pspb.WriteRequest{
		Key:   CONFIG_KEY,
		Value: &anypb.Any{Value: data},
	})
	return err
}

func main() {
	ghclient, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	// Get a pstore client
	pclient, err := pstoreclient.GetClient()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	s := &Server{
		Client:  ghclient,
		Pclient: pclient,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpc_port))
	if err != nil {
		log.Fatalf("gramophile is unable to listen on the grpc port %v: %v", *grpc_port, err)
	}

	// Handle grpc requests
	gs := grpc.NewServer()
	pb.RegisterOverseerServiceServer(gs, s)
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("gramophile is unable to serve grpc: %v", err)
		}
		log.Fatalf("gramophile has closed the grpc port for some reason")
	}()

	go func() {
		s.runBackground()
	}()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf(":%v", *metrics_port), nil)
	log.Fatalf("gramophile is unable to serve metrics: %v", err)

}
