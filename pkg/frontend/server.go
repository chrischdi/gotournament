package frontend

import (
	"net/http"
)

func Run(port string) {
	registerHandlers()

	http.ListenAndServe(port, nil)
}
