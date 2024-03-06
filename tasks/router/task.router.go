package router

import (
	taskcontroller "tasks/controllers"
	authmiddleware "tasks/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authmiddleware.AdminAuthMiddleware)
	adminRouter.HandleFunc("/api/task", taskcontroller.CreateOneTask).Methods("POST")
	adminRouter.HandleFunc("/api/task", taskcontroller.UpdateOne).Methods("PUT")
	adminRouter.HandleFunc("/api/task/{id}", taskcontroller.DeleteOneTask).Methods("DELETE")

	userRouter := router.PathPrefix("/api").Subrouter()
	userRouter.Use(authmiddleware.AuthMiddleware)
	userRouter.HandleFunc("/tasks", taskcontroller.GetMyAllTasks).Methods("GET")

	return router
}
