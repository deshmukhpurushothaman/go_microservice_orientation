package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/deshmukhpurushothaman/golang_microservice_orientation/logger-service/data"
	"github.com/deshmukhpurushothaman/golang_microservice_orientation/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	// return response
	resp := &logs.LogResponse{Result: "logged!"}
	return resp, nil
}

func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRpc server started on port %s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
}
