package main

import (
	"net/http"

	"github.com/coderminer/restful/routes"
)

func main() {
	r := routes.NewRouter()

	http.ListenAndServe(":8080", r)
}
