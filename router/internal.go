package router

// this to route internal requests from admin page
import (
	"github.com/SurgicalSteel/elasthink/service"
	"net/http"
)

//RegisterInternalHandler registers internal handlers (internal endpoints)
func (rw *RouterWrap) RegisterInternalHandler() {
	subRouteInternalV1 := rw.Router.PathPrefix("/internal/v1").Subrouter()

	subRouteInternalV1.HandleFunc("/index/{document_type}/{document_id}", service.HandleCreateIndex).Methods(http.MethodPost)
	subRouteInternalV1.HandleFunc("/index/{document_type}/{document_id}", service.HandleUpdateIndex).Methods(http.MethodPut)

	subRouteInternalV1.HandleFunc("/search/{document_type}", service.HandleSearch).Methods(http.MethodPost)

}
