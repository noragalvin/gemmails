package routes

import (
	"gemmails/app/controllers"

	"github.com/gorilla/mux"
)

// NewRouter function configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// router.Use(auth.ValidateAuth)

	router.HandleFunc("/api/send-mail/{destination}", controllers.MailSend).Methods("POST")

	return router
}
