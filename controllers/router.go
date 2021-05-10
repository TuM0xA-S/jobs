package controllers

import (
	"jobs/util"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//FormatPanicError just responds with json
func (i InternalServerErrorResponder) FormatPanicError(rw http.ResponseWriter, _ *http.Request, _ *negroni.PanicInformation) {
	util.RespondWithError(rw, 500, "server internal error")
}

//InternalServerErrorResponder [lol i hate suppressing that warning messages]
type InternalServerErrorResponder struct{}

// GetRouter returns prepared router
func GetRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/jobs", JobsList).Methods("GET")
	router.HandleFunc("/jobs/{job}", JobDetail).Methods("GET")

	router.HandleFunc("/jobs", JobCreate).Methods("POST")
	router.HandleFunc("/tasks", TaskCreate).Methods("POST").Queries("parent", "{parent}")
	router.HandleFunc("/works", WorkCreate).Methods("POST").Queries("parent", "{parent}")

	router.HandleFunc("/tasks/{task}", TaskDelete).Methods("DELETE")
	router.HandleFunc("/works/{work}", WorkDelete).Methods("DELETE")
	router.HandleFunc("/jobs/{job}", JobDelete).Methods("DELETE")

	router.HandleFunc("/tasks/{task}", TaskMove).Methods("PUT").Queries("parent", "{parent}")
	router.HandleFunc("/works/{work}", WorkMove).Methods("PUT").Queries("parent", "{parent}")

	router.NotFoundHandler = http.HandlerFunc(NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowed)

	n := negroni.New()

	logger := log.New(os.Stdout, "[jobs]", 0)

	loggerMid := negroni.NewLogger()
	loggerMid.ALogger = logger
	n.Use(loggerMid)

	recoverMid := negroni.NewRecovery()
	recoverMid.Logger = logger
	recoverMid.Formatter = InternalServerErrorResponder{}
	n.Use(recoverMid)

	n.UseHandler(router)

	return n
}
