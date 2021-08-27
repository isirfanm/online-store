package main

import (
	"github.com/isirfanm/online-store/api"
	"github.com/isirfanm/online-store/config"
)

func main() {
	// setup all
	config.SetupAll()

	// init router
	router := api.SetupRouter()

	// run gin
	router.Run()
}
