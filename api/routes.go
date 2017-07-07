// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/client",
		http.StripPrefix("/client", http.FileServer(http.Dir("api/safing-client"))),
	},
	Route{
		"Websockets",
		"GET",
		"/api/v1",
		http.HandlerFunc(apiVersionOneHandler),
	},
}
