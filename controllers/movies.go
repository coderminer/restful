package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/coderminer/restful/models"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	dao = models.Movies{}
)

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movies []models.Movies
	movies, err := dao.FindAllMovies()
	if err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, movies)

}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := dao.FindMovieById(id)
	if err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, result)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movies
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.Id = bson.NewObjectId()
	if err := dao.InsertMovie(movie); err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusCreated, movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Movies
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.UpdateMovie(params); err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := dao.RemoveMovie(id); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
