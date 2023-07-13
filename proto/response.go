package proto

import (
	"encoding/json"
	"io"

	"github.com/amar-jay/minicache/errors"
	"github.com/amar-jay/minicache/logger"
)

type Responses interface{}

type Response struct {
	Key    []byte `json:"key,omitempty"`
	Value  []byte `json:"value,omitempty"`
	Status int    `jsoin:"status"` // http status code
}

func (r *Response) Bytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, logger.Errorf(errors.ParseError, err.Error())
	}

	return jsonBytes, nil
}

func ParseSetResponse(r io.Reader) (*Response, error) {
	resp := new(Response)
	if err := json.NewDecoder(r).Decode(resp); err != nil {
		return nil, logger.Errorf(errors.ParseError, err.Error())
	}

	return resp, nil
}
