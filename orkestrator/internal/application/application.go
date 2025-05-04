package application

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/veronicashkarova/server-for-calc/pkg/contract"
	"github.com/veronicashkarova/server-for-calc/pkg/db"
) 


func ConfigFromEnv() *contract.Config {

	config := new(contract.Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	addTime, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err == nil {
		config.TIME_ADDITION_MS = addTime
	} else {
		config.TIME_ADDITION_MS = 1000
	}
	subTime,err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err == nil {
		config.TIME_SUBTRACTION_MS = subTime
	} else {
		config.TIME_SUBTRACTION_MS = 1000
	}
	mulTime, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err == nil {
		config.TIME_MULTIPLICATIONS_MS = mulTime
	} else {
		config.TIME_SUBTRACTION_MS = 1000
	}
	divTime,err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err == nil {
		config.TIME_DIVISIONS_MS = divTime
	} else {
		config.TIME_DIVISIONS_MS = 1000
	}
	return config
}

type Application struct {
	config *contract.Config
}

func New() *Application {
	contract.AppConfig = ConfigFromEnv()
	return &Application{
		config: contract.AppConfig,
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type TaskRequest struct {
	ID int `json:"id"`
	Result float64 `json:"result"`
}

func (a *Application) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", RegisterUserHandler)
	mux.HandleFunc("/api/v1/login", LoginUserHandler)
	mux.HandleFunc("/internal/task", TaskHandler)
	calculate := AutorizationMiddleware(http.HandlerFunc(NewExpressionHandler))
    expressions := AutorizationMiddleware(http.HandlerFunc(ExpressionsHandler))
	idExpressions := AutorizationMiddleware(http.HandlerFunc(IdHandler))
	mux.Handle("/api/v1/calculate", calculate)
	mux.Handle("/api/v1/expressions", expressions)
	mux.Handle("/api/v1/expressions/", idExpressions)
	fmt.Println("Server started")
	return http.ListenAndServe("", mux)
}

func (a *Application) CreareDataBase()  {
	db.CreateDb()
}
