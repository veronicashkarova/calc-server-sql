package application

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/veronicashkarova/server-for-calc/pkg/orkestrator"
	pb "github.com/veronicashkarova/server-for-calc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.CalculatorServiceServer // сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

type CalculatorServiceServer interface {
	GetTask(context.Context, *pb.EmptyRequest) (*pb.Task, error)
	GetResult(context.Context, *pb.TaskResult) *pb.EmptyResponse
	mustEmbedUnimplementedGeometryServiceServer()
}

func (s *Server) GetTask(
	ctx context.Context,
	req *pb.EmptyRequest,
) (*pb.Task, error) {
	task, err := orkestrator.GetTaskData()
	fmt.Println("GetTaskData", task, err)
	if err != nil {
		return &pb.Task{}, err
	}
	return &pb.Task{
		Id:            int32(task.ID),
		Arg1:          float32(task.Arg1),
		Arg2:          float32(task.Arg2),
		Operation:     task.Operation,
		OperationTime: int32(task.OperationTime),
	}, nil
}

func (s *Server) GetResult(
	ctx context.Context,
	taskResult *pb.TaskResult,
) (*pb.EmptyResponse, error) {
	resp := &pb.EmptyResponse{}
	var resultErr = orkestrator.SendResult(int(taskResult.Id), float64(taskResult.Result))
	if resultErr != nil {
		return resp, resultErr
	}
	return resp, nil
}

func StartGrpcServer() {
	go func() {
		host := "localhost"
		port := "5000"

		addr := fmt.Sprintf("%s:%s", host, port)
		lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

		if err != nil {
			fmt.Println("error starting tcp listener: ", err)
			os.Exit(1)
		}

		fmt.Println("tcp listener started at port: ", port)
		// создадим сервер grpc
		grpcServer := grpc.NewServer()
		// объект структуры, которая содержит реализацию
		// серверной части GeometryService
		calcServiceServer := NewServer()
		// зарегистрируем нашу реализацию сервера
		pb.RegisterCalculatorServiceServer(grpcServer, calcServiceServer)
		// запустим grpc сервер
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("error serving grpc: ", err)
			os.Exit(1)
		}
	}()
}
