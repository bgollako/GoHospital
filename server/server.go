package server

import (
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/db"
	"context"
	"encoding/json"
	"net/http"
)

type CustomHandler struct {
}

func (s *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.RequestURI == "/v1/Patients" {
			patients, err := db.GetAllPatients(context.Background())
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Internal Error " + err.Error()))
				return
			}
			data, err := json.Marshal(patients)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Internal Error " + err.Error()))
				return
			}
			w.WriteHeader(200)
			w.Write(data)
			return
		}

		w.WriteHeader(400)
		w.Write([]byte("Invalid URL"))
		return
	}
	w.WriteHeader(405)
	w.Write([]byte("Unsupported method"))
	return
}
