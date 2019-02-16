package main

import (
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/db"
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/server"
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	handler := &server.CustomHandler{}
	httpServer := http.Server{
		Addr:         ":8082",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Starting Server ...")
	err := db.TestConnection(context.Background())
	if err != nil {
		fmt.Println("Error while pinging mongo db " + err.Error())
		return
	}
	fmt.Println("Successfully pinged mongodb")
	err = db.InitialiseDummyPatients(context.Background())
	if err != nil {
		fmt.Println("Error while initializing patients " + err.Error())
		return
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error while starting http server : " + err.Error())
	}
}
