package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func EnrollmentResponseIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.EnrollmentResponse
	c := Database.C("enrollmentresponses")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var enrollmentresponseEntryList []models.BundleEntryComponent
	for _, enrollmentresponse := range result {
		var entry models.BundleEntryComponent
		entry.Resource = &enrollmentresponse
		enrollmentresponseEntryList = append(enrollmentresponseEntryList, entry)
	}

	var bundle models.Bundle
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Type = "searchset"
	var total = uint32(len(result))
	bundle.Total = &total
	bundle.Entry = enrollmentresponseEntryList

	log.Println("Setting enrollmentresponse search context")
	context.Set(r, "EnrollmentResponse", result)
	context.Set(r, "Resource", "EnrollmentResponse")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(&bundle)
}

func LoadEnrollmentResponse(r *http.Request) (*models.EnrollmentResponse, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("enrollmentresponses")
	result := models.EnrollmentResponse{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting enrollmentresponse read context")
	context.Set(r, "EnrollmentResponse", result)
	context.Set(r, "Resource", "EnrollmentResponse")
	return &result, nil
}

func EnrollmentResponseShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadEnrollmentResponse(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "EnrollmentResponse"))
}

func EnrollmentResponseCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	enrollmentresponse := &models.EnrollmentResponse{}
	err := decoder.Decode(enrollmentresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("enrollmentresponses")
	i := bson.NewObjectId()
	enrollmentresponse.Id = i.Hex()
	err = c.Insert(enrollmentresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting enrollmentresponse create context")
	context.Set(r, "EnrollmentResponse", enrollmentresponse)
	context.Set(r, "Resource", "EnrollmentResponse")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/EnrollmentResponse/"+i.Hex())
}

func EnrollmentResponseUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	enrollmentresponse := &models.EnrollmentResponse{}
	err := decoder.Decode(enrollmentresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("enrollmentresponses")
	enrollmentresponse.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, enrollmentresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting enrollmentresponse update context")
	context.Set(r, "EnrollmentResponse", enrollmentresponse)
	context.Set(r, "Resource", "EnrollmentResponse")
	context.Set(r, "Action", "update")
}

func EnrollmentResponseDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("enrollmentresponses")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting enrollmentresponse delete context")
	context.Set(r, "EnrollmentResponse", id.Hex())
	context.Set(r, "Resource", "EnrollmentResponse")
	context.Set(r, "Action", "delete")
}
