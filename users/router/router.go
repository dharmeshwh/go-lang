package router

import (
	usercontroller "users/controller"
	authmiddleware "users/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", usercontroller.Login).Methods("POST")

	adminRouter := router.PathPrefix("/api").Subrouter()
	adminRouter.Use(authmiddleware.AdminAuthMiddleware)
	adminRouter.HandleFunc("/users", usercontroller.GetMyAllUsers).Methods("GET")
	adminRouter.HandleFunc("/user", usercontroller.CreateUser).Methods("POST")
	adminRouter.HandleFunc("/user/{id}", usercontroller.DeleteAUser).Methods("DELETE")

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(authmiddleware.AuthMiddleware)
	userRouter.HandleFunc("/profile", usercontroller.GetUserDetails).Methods("GET")

	return router
}
