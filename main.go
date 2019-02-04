package main

import (
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/server"
	"fmt"
	"net/http"
	"time"
)

func main() {
	handler := &server.CustomHandler{}
	httpServer := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Starting Server ...")
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error while starting http server : " + err.Error())
	}
}
