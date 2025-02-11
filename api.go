package main

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/brotherlogic/overseer/proto"
)

func (s *Server) RegisterTask(ctx context.Context, req *pb.RegisterTaskRequest) (*pb.RegisterTaskResponse, error) {
	config, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	uuid := uuid.New().String()

	config.Tasks = append(config.Tasks, &pb.OverseerTask{
		Uuid:     uuid,
		Callback: req.GetCallback(),
		Task:     req.GetTask(),
	})

	return &pb.RegisterTaskResponse{}, s.saveConfig(ctx, config)
}
