package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const auth = "0031RBsvIzZifBRtSKJ2dJ9qZMVXgccQg0G4c23YvpsH1H"
const service = "http://icfpc2013.cloudapp.net/"

type TrainRequest struct {
	Size      int      `json:"size"`
	Operators []string `json:"operators"`
}

type TrainingProblem struct {
	Challenge string
	Id        string
	Size      int
	Operators []string
}

type EvalRequest struct {
	Id        string   `json:"id"`
	Program   string   `json:"program"`
	Arguments []string `json:"arguments"`
}

type EvalResponse struct {
	Status  string
	Outputs []string
	Message string
}

type Problem struct {
    Id string
    Size int
    Operators []string
    Solved bool
    TimeLeft int
}

func request(kind string, body []byte) ([]byte, error) {
	resp, err := http.Post(service+kind+"?auth="+auth, "text/plain", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rbody, nil
}

func Train(req TrainRequest) (*TrainingProblem, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := request("train", body)
	if err != nil {
		return nil, err
	}
	var tp TrainingProblem
	if err := json.Unmarshal(resp, &tp); err != nil {
    return nil, err
  } else {
    return &tp, nil
  }
}

func Eval(req EvalRequest) (*EvalResponse, error) {
  body, err := json.Marshal(req)
  if err != nil {
    return nil, err
  }
  resp, err := request("eval", body)
  if err != nil {
    return nil, err
  }
  var er EvalResponse
  if err := json.Unmarshal(resp, &er); err != nil {
    return nil, err
  } else {
    return &er, nil
  }
}
