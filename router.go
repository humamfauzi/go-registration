package main

import "encoding/json"

type OperationReply struct {
	Name    string
	Success bool
}

func (or *OperationReply) SetFail() {
	or.Success = false
}

type Reply struct {
	Operation OperationReply `json:"operation"`
	Error     ErrorReply     `json:"error"`
	Body      interface{}    `json:"body,omitempty"`
}

type ErrorReply struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    string `json:"meta,omitempty"`
}

func CreateReply(opProfile OperationReply, erProfile ErrorReply, body interface{}) ([]byte, error) {
	newReply := Reply{
		Operation: opProfile,
		Error:     erProfile,
		Body:      body,
	}
	reply, err := json.Marshal(newReply)
	if err != nil {
		return []byte{}, err
	}
	return reply, nil
}
