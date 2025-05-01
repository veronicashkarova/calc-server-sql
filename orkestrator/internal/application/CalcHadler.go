package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/veronicashkarova/server-for-calc/pkg/calc"
	"github.com/veronicashkarova/server-for-calc/pkg/contract"
	"github.com/veronicashkarova/server-for-calc/pkg/orkestrator"
)

func NewExpressionHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, id, err := orkestrator.AddExpression(request.Expression)
	if err != nil {
		switch {
		case errors.Is(err, calc.ErrInvalidExpression):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, calc.ErrEmptyExpression):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, result)

		go func() {
			result, err := calc.Calc(request.Expression, id, contract.TaskChannel)
			value, exist := contract.ExpressionMap[id]
			if exist {
				if err != nil {
					value.Data.Status = err.Error()
				} else {
					value.Data.Status = contract.Done
					value.Data.Result = strconv.FormatFloat(result, 'f', 3, 64)
				}
				contract.ExpressionMap[id] = value
			}
		}()
	}
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	result, err := orkestrator.Expressions()
	if err != nil {
		switch {
		case errors.Is(err, calc.ErrInvalidExpression):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, calc.ErrEmptyExpression):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result)
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := isIdExpressionRequest(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	result, err := orkestrator.GetExpressionForId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		result, err := orkestrator.GetTask()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Fprint(w, result)

	} else {
		taskRequest := new(TaskRequest)
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&taskRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resultErr = orkestrator.SendResult(taskRequest.ID, taskRequest.Result)
		if resultErr != nil {
			http.Error(w, resultErr.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func isIdExpressionRequest(url *url.URL) (string, error) {

	// Разделяем путь на сегменты
	pathSegments := strings.Split(url.Path, "/")

	// Получаем последний сегмент
	lastSegment := pathSegments[len(pathSegments)-1]
	fmt.Println(lastSegment)

	return lastSegment, nil
}
