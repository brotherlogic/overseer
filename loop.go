package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/overseer/proto"
)

const (
	DELAY = time.Minute * 5
)

func (s *Server) runBackground() {
	for {
		time.Sleep(DELAY)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		err := s.runLoop(ctx)
		if err != nil {
			log.Printf("Loop returned error: %v", err)
		}
		cancel()
	}
}

// Runs a loop of overseer - finds a single validation tasks if possible and
// runs that on a schedule
func (s *Server) runLoop(ctx context.Context) error {
	config, err := s.loadConfig(ctx)
	if err != nil {
		return err
	}

	for _, task := range config.GetTasks() {
		lastRun := int64(0)
		for _, validation := range task.GetValidationRuns() {
			if validation.GetTimestampMs() > lastRun {
				lastRun = validation.GetTimestampMs()
			}
		}

		if time.Since(time.UnixMilli(lastRun)).Seconds() > float64(task.GetDelayInS()) {
			return s.runTask(ctx, task, config)
		}
	}

	return nil
}

func (s *Server) runTask(ctx context.Context, task *pb.OverseerTask, config *pb.Config) error {
	conn, err := grpc.NewClient(task.GetCallback(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewOverseerClientClient(conn)
	res, err := client.Overseer(ctx, &pb.OverseerRequest{
		Task: task.GetTask(),
	})

	if err != nil {
		task.ValidationRuns = append(task.ValidationRuns, &pb.ValidationRun{
			TimestampMs:   time.Now().UnixMilli(),
			CanonicalCode: fmt.Sprintf("%v", status.Code(err)),
		})
	} else {
		task.ValidationRuns = append(task.ValidationRuns, &pb.ValidationRun{
			TimestampMs: time.Now().UnixMilli(),
			Response:    res.GetResponse(),
		})
	}

	return s.saveConfig(ctx, config)
}
