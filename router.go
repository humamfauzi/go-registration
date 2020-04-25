package main

type OperationReply struct {
	Name    string
	Success bool
}

func(or *OperationReply) Flip() {
	if or.Success {
		or.Success = false
	} else {
		or.Success = true
	}
}
type Reply struct {
	Operation OperationReply
	Error     ErrorReply
	Body      interface{}
}

type ErrorReply struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    string `json:"meta,omitempty"`
}


func CreateReply(opProfile OperationReply, erProfile ErrorReply, body interface{}) ([]byte, error) {
	newReply := Reply{
		Operation: opProfile,
		Error: erProfile,
		Body: body
	}
	reply, err := json.Marshal(newReply)
	if err != nil {
		return []byte{}, err
	}
	return reply, nil
}