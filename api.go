package icfp13

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const auth = "0031RBsvIzZifBRtSKJ2dJ9qZMVXgccQg0G4c23YvpsH1H"
const service = "http://icfpc2013.cloudapp.net/"

type TrainingProblem struct {
	Challenge string
	Id        string
	Size      int
	Operators []string
}

type EvalResponse struct {
	status  string
	outputs []string
	message string
}

func request(kind, body string) ([]byte, error) {
	resp, err := http.Post(service+kind+"?auth="+auth, "text/plain", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rbody, nil
}

func Train(size int, operators []string) (*TrainingProblem, error) {
	resp, err := request("train", "{ \"size\": 3, \"operators\": [] }")
	if err != nil {
		return nil, err
	}
	fmt.Println("response:", string(resp))
	var tp TrainingProblem
	json.Unmarshal(resp, &tp)
	return &tp, nil
}
