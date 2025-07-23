package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case http.MethodPost:
		addTaskHandler(w, req)

	case http.MethodGet:
		getTaskHandler(w, req)

	case http.MethodPut:
		updTaskHandler(w, req)

	case http.MethodDelete:
		delTaskHandler(w, req)
	}
}
