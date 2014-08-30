package HelloService

import (
	"log"
	"net/http"
)

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {

	reply.Message = "Hello, " + args.Who + "!"
	log.Printf(reply.Message)
	return nil
}

func (h *HelloService) ParameterReflection(r *http.Request, args *interface{}, reply *interface{}) error {

	// Create the return object
	ret := make(map[string]interface{})
	ret["reflection"] = args
	*reply = &ret

	return nil
}
