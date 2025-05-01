package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
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

func RunAgent(power int, delay int) {
	fmt.Println("start agent")

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startAgent(delay)
		}()
	}

	wg.Wait()
}

func startAgent(delay int) {
	orchestratorURL := "http://localhost"

	for {
		task, err := getTask(orchestratorURL + "/internal/task")
		if err != nil {
			log.Printf("Ошибка получения задачи. Повторная попытка через %d секунд...", delay/1000)
			Delay(delay)
			continue
		}

		operationTimer := time.NewTimer(time.Duration(task.OperationTime * int(time.Millisecond)))
		<-operationTimer.C

		result, err := executeTask(task)
		if err != nil {
			log.Printf("Ошибка выполнения задачи")
			continue
		}

		err = submitResult(orchestratorURL+"/internal/task", &result)
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

func getTask(url string) (Task, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Task{}, fmt.Errorf("неуспешный код ответа: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Task{}, err
	}

	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return Task{}, err
	}

	fmt.Println("task", task)

	return task, nil
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

func submitResult(url string, result *Result) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
