package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/globalsign/mgo/bson"

	"github.com/coderminer/restful/models"
)

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implmented !")
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movies
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.Id = bson.NewObjectId()
	if err := movie.InsertMovie(movie); err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusCreated, movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}
