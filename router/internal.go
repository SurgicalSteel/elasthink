package router

// Elasthink, An alternative to elasticsearch engine written in Go for small set of documents that uses inverted index to build the index and utilizes redis to store the indexes.
// Copyright (C) 2020 Yuwono Bangun Nagoro (a.k.a SurgicalSteel)
//
// This file is part of Elasthink
//
// Elasthink is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Elasthink is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

}
