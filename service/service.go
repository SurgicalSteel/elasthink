//Package service is where we handle every incoming services that are defined in router package
package service

import (
	"github.com/SurgicalSteel/elasthink/module"
	"net/http"
)

//HandlePing is the handler for a ping endpoint
func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG!"))
}

//ResponsePayload is the response payload struct for all API
type ResponsePayload struct {
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

func constructResponsePayload(rawModuleResponse module.Response) ResponsePayload {
	return ResponsePayload{
		ErrorMessage: rawModuleResponse.ErrorMessage,
		Data:         rawModuleResponse.Data,
	}
}
