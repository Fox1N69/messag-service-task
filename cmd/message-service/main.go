package main

import (
	"messaggio/infra"
	"messaggio/internal/api"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	// Init config
	i := infra.New("config/config.json")
	// Set project mod
	i.SetMode()

	// Get custom logrus logger
	log := i.GetLogger()

	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//Connect to database and migration
	i.PSQLClient()
	log.Info("Connect to PSQLClient")

	// Start api server
	api.NewServer(i).Run()
}
