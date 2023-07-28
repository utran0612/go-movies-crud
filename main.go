package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Ibsn     string    `json:"ibsn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

// this is a HTTP request handler
func getMovies(w http.ResponseWriter, r *http.Request) {

	// w is to write the response back to client (the one that requests)
	w.Header().Set("Content-Type", "application/json")

	//encode w into json
	//json.NewEncoder() creates a NewEncoder()
	//then use Encode method to encode it into json
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the params from the request
	params := mux.Vars(r)
	for i, item := range movies {
		if item.ID == params["id"] {
			//take the array of movie 1 -> i-1
			//and append to the array of movie i+1 -> n
			//so movie[i] doesn't exist in movies anymore
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
	return
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decode movie from json to the struct format
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies, movie)

	//encode the movie from our defined struct to json again
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// Search for the movie with the matching ID in the global slice
	for index, item := range movies {
		if item.ID == movie.ID {
			// Update the movie's title and/or director if they are different
			if item.Title != movie.Title {
				movies[index].Title = movie.Title
			}
			if item.Director != movie.Director {
				movies[index].Director = movie.Director
			}
			break
		}
	}

	// Return the updated movie as the response
	json.NewEncoder(w).Encode(movie)
}

func main() {
	r := mux.NewRouter()

	//add some movies
	movies = append(movies, Movie{ID: "1", Ibsn: "438227", Title: "Barbie", Director: &Director{FirstName: "John", LastName: "Barh"}})
	movies = append(movies, Movie{ID: "2", Ibsn: "45455", Title: "Openheimmer", Director: &Director{FirstName: "Leila", LastName: "Kim"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server on port 8000\n")

	//starts an HTTP server with a given address and handler
	//this case address is port 8000, handler is r (which uses gorilla mux)
	//if handler is nil, then it uses DefaultServerMux
	log.Fatal(http.ListenAndServe(":8000", r)) //ListenAndServe returns an error
}
