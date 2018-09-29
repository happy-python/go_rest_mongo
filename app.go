package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	. "go_rest_mongo/config"
	. "go_rest_mongo/dao"
	. "go_rest_mongo/models"
	utils "go_rest_mongo/utils"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var config = Config{}
var dao = MoviesDAO{}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// GET list of movies
func AllMoviesEndPoint(w http.ResponseWriter, req *http.Request) {
	movies, err := dao.FindAll()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusCreated, movies)
}

// GET a movie by its ID
func FindMovieEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	movie, err := dao.FindById(params["id"])

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, movie)

}

// POST a new movie
func CreateMovieEndpoint(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var movie Movie
	if err := json.NewDecoder(req.Body).Decode(&movie); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusCreated, movie)
}

// PUT update an existing movie
func UpdateMovieEndPoint(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var movie Movie
	if err := json.NewDecoder(req.Body).Decode(&movie); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.Update(movie); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing movie
func DeleteMovieEndPoint(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var movie Movie
	if err := json.NewDecoder(req.Body).Decode(&movie); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.Delete(movie); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndpoint).Methods("POST")
	r.HandleFunc("/movies/{id}", FindMovieEndpoint).Methods("GET")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
