package json_multipart

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/rpc"
	"net/http"
)

var null = json.RawMessage([]byte("null"))

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

type serverResponse struct {
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
	Id     *json.RawMessage `json:"id"`
}

// ----------------------------------------------------------------------------
// Codec
// ----------------------------------------------------------------------------

// NewCodec returns a new JSON Codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Codec creates a CodecRequest to process each request.
type Codec struct {
}

// NewRequest returns a CodecRequest.
func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	return newCodecRequest(r)
}

// ----------------------------------------------------------------------------
// CodecRequest
// ----------------------------------------------------------------------------

// newCodecRequest returns a new CodecRequest.
func newCodecRequest(r *http.Request) rpc.CodecRequest {
	// Decode the request body and check if RPC method is valid.
	req := new(serverRequest)

	req.Method = r.FormValue("method")

	paramsValue := json.RawMessage(r.FormValue("params"))
	req.Params = &paramsValue

	idval := json.RawMessage(r.FormValue("id"))
	req.Id = &idval

	var err error = nil
	//	err := json.NewDecoder(r.Body).Decode(req)
	r.Body.Close()
	return &CodecRequest{request: req, err: err}
}

// CodecRequest decodes and encodes a single request.
type CodecRequest struct {
	request *serverRequest
	err     error
}

// Method returns the RPC method for the current request.
//
// The method uses a dotted notation as in "Service.Method".
func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

// ReadRequest fills the request object for the RPC method.
func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.err == nil {
		if c.request.Params != nil {
			// JSON params is array value. RPC params is struct.
			// Unmarshal into array containing the request struct.
			params := [1]interface{}{args}
			c.err = json.Unmarshal(*c.request.Params, &params)
		} else {
			c.err = errors.New("rpc: method request ill-formed: missing params field")
		}
	}
	return c.err
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
//
// The err parameter is the error resulted from calling the RPC method,
// or nil if there was no error.
func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, methodErr error) error {
	if c.err != nil {
		return c.err
	}
	res := &serverResponse{
		Result: reply,
		Error:  &null,
		Id:     c.request.Id,
	}
	if methodErr != nil {
		// Propagate error message as string.
		res.Error = methodErr.Error()
		// Result must be null if there was an error invoking the method.
		// http://json-rpc.org/wiki/specification#a1.2Response
		res.Result = &null
	}
	if c.request.Id == nil {
		// Id is null for notifications and they don't have a response.
		res.Id = &null
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}
	return nil
}
