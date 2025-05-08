package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "github.com/veronicashkarova/agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Result struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

func RunGrpcAgent(power int, delay int) {
	fmt.Println("start agent")

	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := pb.NewCalculatorServiceClient(conn)

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startGrpcAgent(grpcClient, delay)
		}()
	}

	wg.Wait()
}

func startGrpcAgent(client pb.CalculatorServiceClient, delay int) {
	ctx := context.TODO()

	for {
		req, err := client.GetTask(ctx, &pb.EmptyRequest{})

		if err != nil {
			log.Printf("Ошибка получения задачи. Повторная попытка через %d секунд...", delay/1000)
			Delay(delay)
			continue
		}

		task := Task{
			ID:            int(req.Id),
			Arg1:          float64(req.Arg1),
			Arg2:          float64(req.Arg2),
			Operation:     req.Operation,
			OperationTime: int(req.OperationTime),
		}

		operationTimer := time.NewTimer(time.Duration(task.OperationTime * int(time.Millisecond)))
		<-operationTimer.C

		result, err := executeTask(task)
		if err != nil {
			log.Printf("Ошибка выполнения задачи")
			continue
		}

		_, err = client.GetResult(ctx, &pb.TaskResult{
			Id:     int32(result.ID),
			Result: float32(result.Result),
		})

		if err != nil {
			log.Printf("Ошибка отправки результата задачи")
			Delay(delay)
			continue
		}

		fmt.Printf("Задача %d выполнена успешно. Результат: %f\n", task.ID, result.Result)
	}
}

func Delay(delay int) {
	delayTimer := time.NewTimer(time.Duration(delay * int(time.Millisecond)))
	<-delayTimer.C
}

func executeTask(task Task) (Result, error) {
	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		result = task.Arg1 / task.Arg2
	default:
		return Result{}, fmt.Errorf("неизвестная операция: %s", task.Operation)
	}

	return Result{ID: task.ID, Result: result}, nil
}
