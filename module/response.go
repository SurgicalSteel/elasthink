package module

import ()

//Response is the universal response struct for all API
type Response struct {
	StatusCode   int
	ErrorMessage string
	Data         interface{}
}
