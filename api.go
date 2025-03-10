package main

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/overseer/proto"
)

func (s *Server) RegisterTask(ctx context.Context, req *pb.RegisterTaskRequest) (*pb.RegisterTaskResponse, error) {
	config, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Don't append if we already have this task
	for _, task := range config.GetTasks() {
		if task.GetTask() == req.GetTask() {
			return nil, status.Errorf(codes.AlreadyExists, "This task (%v) already exists", req.GetTask())
		}
	}

	uuid := uuid.New().String()

	config.Tasks = append(config.Tasks, &pb.OverseerTask{
		Uuid:     uuid,
		Callback: req.GetCallback(),
		Task:     req.GetTask(),
	})

	return &pb.RegisterTaskResponse{}, s.saveConfig(ctx, config)
}
