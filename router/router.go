//package router is the package where we define web routes on elasthink
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
import (
	"github.com/SurgicalSteel/elasthink/service"
	"github.com/gorilla/mux"
	"net/http"
)

//RouterWrap is a custom wrapper type for router, you can add more configuration fields here
type RouterWrap struct {
	Router *mux.Router
}

var routeWrap *RouterWrap

// RegisterHandler is a RouterWrap 'method' to register your API endpoints.
// Usually handler calls services module
func (rw *RouterWrap) RegisterHandler() {
	rw.Router.HandleFunc("/ping", service.HandlePing).Methods(http.MethodGet)
}

// InitializeRoute is a function which returns new RouterWrap which has a mux's router inside
func InitializeRoute() *RouterWrap {
	routeWrap = new(RouterWrap)
	routeWrap.Router = mux.NewRouter()
	return routeWrap
}
