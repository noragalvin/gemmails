package routes

import (
	"gemmails/app/controllers"

	"github.com/gorilla/mux"
)

// NewRouter function configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// router.Use(auth.ValidateAuth)

	/* Theme router
	--Default API router must be:
	---/api/themes -> GET: get list themes
	---/api/themes/{id} -> GET: get a theme by id
	---/api/themes -> POST: store new theme
	---/api/themes/{id} -> PUT: edit theme by id
	---/api/themes/{id} -> DELETE: destroy theme by id
	*/
	// router.HandleFunc("/api/admin/theme", controllers.ThemeAdminUpdate).Methods("POST")

	router.HandleFunc("/api/send-mail/{destination}", controllers.MailSend).Methods("GET")

	return router
}
