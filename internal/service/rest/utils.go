package rest

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (er ErrorResponse) Raw() []byte {
	packed, _ := json.Marshal(er)
	return packed
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Status: "err",
		Reason: err.Error(),
	}
}

func response(resp any, code int, w http.ResponseWriter) {
	var out []byte

	out, err := json.Marshal(resp)
	if err != nil {
		out = NewErrorResponse(err).Raw()
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	w.Write(out)
}
