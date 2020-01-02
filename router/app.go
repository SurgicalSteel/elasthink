package router

import (
	"github.com/SurgicalSteel/elasthink/service"
	"net/http"
)

//RegisterAppHandler register app handlers (external endpoints)
func (rw *RouterWrap) RegisterAppHandler() {
	subRouteV1 := rw.Router.PathPrefix("/v1").Subrouter()

	subRouteV1.HandleFunc("/search/{document_type}", service.HandleSearch).Methods(http.MethodPost)
}
