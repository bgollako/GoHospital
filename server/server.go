package server

import (
	"bitbucket-eng-sjc1.cisco.com/an/GoHospital/db"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type CustomHandler struct {
}

func InvalidUrl(w http.ResponseWriter)  {
	w.WriteHeader(400)
	w.Write([]byte("Invalid URL"))
}

func InternalError(w http.ResponseWriter, err error)  {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func Success(w http.ResponseWriter, data []byte)  {
	w.WriteHeader(200)
	w.Write(data)
}

func (s *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodPost:
		if r.RequestURI != "/v1/Patients" {
			InvalidUrl(w)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			InternalError(w, err)
			return
		}
		patients := make([]interface{}, 0)
		err = json.Unmarshal(data, &patients)
		if err != nil {
			InternalError(w, err)
			return
		}
		if len(patients) == 0 {
			InternalError(w, errors.New("no patients specified"))
			return
		}
		err = db.PostPatient(ctx, patients)
		if err != nil{
			InternalError(w, err)
			return
		}
		Success(w, []byte("Patient(s) creation successful"))
	case http.MethodGet:
		if r.RequestURI != "/v1/Patients" {
			InvalidUrl(w)
			return
		}
		patients, err := db.GetAllPatients(ctx)
		if err != nil {
			InternalError(w, err)
			return
		}
		data, err := json.Marshal(patients)
		if err != nil {
			InternalError(w, err)
			return
		}
		Success(w, data)
	case http.MethodDelete:
		if !strings.HasPrefix(r.RequestURI, "/v1/Patients") {
			InvalidUrl(w)
			return
		}
		if !strings.HasPrefix(r.URL.RawQuery, "id") {
			InvalidUrl(w)
			return
		}
		parts := strings.Split(r.URL.RawQuery, "=")
		if len(parts) != 2 {
			InvalidUrl(w)
			return
		}
		count, err := db.DeletePatient(ctx, parts[1])
		if err != nil {
			InternalError(w, err)
			return
		}
		if count == 0 {
			InternalError(w, errors.New("no patients deleted"))
			return
		}
		Success(w, []byte("Patient(s) deleted successfully"))
	case http.MethodPatch:
		if !strings.HasPrefix(r.RequestURI, "/v1/Patients") {
			InvalidUrl(w)
			return
		}
		if !strings.HasPrefix(r.URL.RawQuery, "id") {
			InvalidUrl(w)
			return
		}
		parts := strings.Split(r.URL.RawQuery, "=")
		if len(parts) != 2 {
			InvalidUrl(w)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			InternalError(w, err)
			return
		}
		patient := &db.Patient{}
		err = json.Unmarshal(data, patient)
		if err != nil {
			InternalError(w, err)
			return
		}
		count, err := db.UpdatePatient(ctx, parts[1], patient)
		if err != nil {
			InternalError(w, err)
			return
		}
		if count == 0 {
			InternalError(w, errors.New("no patients modified"))
			return
		}
		Success(w, []byte("Patient updated successfully"))
	default:
		w.WriteHeader(405)
		w.Write([]byte("Unsupported method"))
	}
}
