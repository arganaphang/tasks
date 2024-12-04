package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "github.com/arganaphang/tasks/gen_proto/task"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	nanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type application struct {
	pb.UnimplementedTaskServiceServer
	tasks []*pb.Task
}

func main() {
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("failed to listen on port 8001: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterTaskServiceServer(s, &application{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"0.0.0.0:8001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = pb.RegisterTaskServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8000",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8000")
	log.Fatalln(gwServer.ListenAndServe())
}

func (app *application) Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	task := &pb.Task{
		Id:          nanoid.Must(),
		Description: req.Description,
		UserId:      req.UserId,
		Status:      pb.TaskStatus_TASK_STATUS_CREATED,
		CreatedAt:   timestamppb.Now(),
	}
	app.tasks = append(app.tasks, task)
	return task, nil
}

func (app *application) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	for idx := range app.tasks {
		if app.tasks[idx].Id == req.Id {
			return app.tasks[idx], nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "task not found")
}

func (app *application) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	return &pb.ListTasksResponse{
		Tasks: app.tasks,
	}, nil
}
