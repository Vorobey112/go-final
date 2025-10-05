package server

import (
	"fmt"
	"net/http"

	"github.com/Vorobey112/go-final/pkg/api"
)

func Run() error {
	port := 7540
	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))
	fmt.Printf("Starting server on http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
