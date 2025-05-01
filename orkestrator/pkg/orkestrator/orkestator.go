package orkestrator

import (
	"encoding/json"
	"fmt"
	"github.com/veronicashkarova/server-for-calc/pkg/calc"
	"github.com/veronicashkarova/server-for-calc/pkg/contract"
	"strconv"
)

var id = 0
var expressions = contract.ExpressionsData{}

func AddExpression(expression string) (string, string, error) {
	id++
	newId := strconv.Itoa(id)
	expressionData :=
		contract.ExpressionData{
			ID:     newId,
			Status: contract.InProcess,
			Result: contract.Undefined,
		}

	contract.ExpressionMap[newId] = contract.ExpressionMapData{
		Data:    expressionData,
		ExpChan: make(chan float64),
	}

	response := contract.ResponseData{ID: newId}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return string(jsonBytes), newId, nil

}

func Expressions() (string, error) {

	for _, expression := range contract.ExpressionMap {
		expressions.Expressions = append(expressions.Expressions, expression.Data)
	}

	jsonBytes, err := json.Marshal(expressions)
	if err != nil {
		panic(err)
	}

	return string(jsonBytes), nil
}

func GetExpressionForId(id string) (string, error) {

	expression, error := findExpressionForId(id)

	if error == nil {
		jsonBytes, err := json.Marshal(expression)
		if err != nil {
			panic(err)
		}
		return string(jsonBytes), nil
	}
	return "", error
}

func GetTask() (string, error) {
	select {
	case task := <-contract.TaskChannel:
		jsonBytes, err := json.Marshal(task)
		if err != nil {
			panic(err)
		}
	
		return string(jsonBytes), nil
	default:
	 return "", calc.ErrNotTask
	}
}

func findExpressionForId(id string) (contract.ExpressionData, error) {
	name, found := contract.ExpressionMap[id]
	if !found {
		return contract.ExpressionData{}, calc.ErrNotFound
	}

	return name.Data, nil
}

func SendResult(id int, result float64) error {
	_, exists := contract.ExpressionMap[fmt.Sprint(id)]
	if exists {
		task := contract.ExpressionMap[fmt.Sprint(id)].Data
		if (task.Result != contract.Done) {
			contract.ExpressionMap[fmt.Sprint(id)].ExpChan <- result
			return nil
		} 
	}
	return calc.ErrNotFound
}
