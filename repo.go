package main

import (
	"context"
	"log"

	pb "github.com/brotherlogic/overseer/proto"
)

func getLatestRun(task *pb.OverseerTask) pb.ValidationResponse {
	first := int64(0)
	var resp pb.ValidationResponse
	for _, run := range task.GetValidationRuns() {
		if run.GetTimestampMs() > first {
			first = run.GetTimestampMs()
			resp = run.GetResponse()
		}
	}

	return resp
}

func (s *Server) runRepoLoop(ctx context.Context, config *pb.Config) error {
	// Find a repo that needs an update
	repos := make(map[string]bool)
	for _, task := range config.GetTasks() {
		repos[task.GetRepo()] = true
	}

	for repo, _ := range repos {
		foundFailure := false
		for _, task := range config.GetTasks() {
			if task.GetRepo() == repo {
				latest := getLatestRun(task)

				if latest != pb.ValidationResponse_VALIDATION_PASSED {
					foundFailure = true
				}
			}
		}

		if !foundFailure {
			log.Printf("Found updateable repo in %v", repo)
		}
	}

	return nil
}
