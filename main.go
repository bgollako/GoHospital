package main

import (
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/db"
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/server"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

func Start()  {
	handler := &server.CustomHandler{}
	httpServer := http.Server{
		Addr:         ":8082",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	ctx := context.Background()
	fmt.Println("Starting Server ...")
	err := db.TestConnection(ctx)
	if err != nil {
		fmt.Println("Error while pinging mongo db " + err.Error())
		return
	}
	defer db.ShutDownClient(ctx)
	fmt.Println("Successfully pinged mongodb")
	//err = db.InitialiseDummyPatients(ctx)
	//if err != nil {
	//	fmt.Println("Error while initializing patients " + err.Error())
	//	return
	//}
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error while starting http server : " + err.Error())
	}
}

func main() {
	Start()
}

func reflectOnStruct()  {
	p := db.Patient{}
	t := reflect.TypeOf(p)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Name)
	}
}



