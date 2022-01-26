package wasm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.dedis.ch/dela/core/store"
	"go.dedis.ch/dela/core/execution"
)

type WASMService struct{}

func (s *WASMService) Execute(snap store.Snapshot, step execution.Step) (execution.Result, error) {
	responseBody := bytes.NewBuffer(step.Current.GetArg("json"))
	resp, err := http.Post("http://127.0.0.1:3000/", "application/json", responseBody)
	if err != nil {
		return execution.Result{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Default()
	args := make(map[string]interface{})
	json.Unmarshal(body, &args)
	acceptedType := fmt.Sprintf("%T", args["Accepted"])
	if acceptedType != "string" {
		return execution.Result{}, errors.New("The value of \"Accepted\" is empty or of a wrong type")
	}
	resultType := fmt.Sprintf("%T", args["result"])
	message := ""
	if resultType == "string" {
		message = args["Accepted"].(string)
	}
	if message == "true" {
		return execution.Result{true, args["result"].(string)}, nil
	}
	return execution.Result{false, args["result"].(string)}, nil
}