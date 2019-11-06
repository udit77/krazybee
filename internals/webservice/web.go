package webservice

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/controllers"
	"github.com/krazybee/internals/dbaccess"
	"log"
	"net/http"
)

type Param struct {
	ConfProvider  config.Provider
	DBProvider dbaccess.Provider
}

func Run(param *Param) error{
	if param == nil {
		return errors.New("[WebService] server expects a valid RouterParam")
	}
	if param.ConfProvider == nil {
		return errors.New("[WebService] server expects a valid configuration")
	}
	if param.DBProvider == nil {
		return errors.New("[WebService] server expects a valid stats pusher")
	}
	searchController, err := controllers.NewController(param.ConfProvider, param.DBProvider)
	if err != nil{
		return errors.New("[WebService] error initializing search controller")
	}
	log.Println("[WebService] starting search service app")
	router := mux.NewRouter()
	router.HandleFunc("/search", searchController.Search).Methods("GET")
	server := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:9000",
	}
	log.Fatal(server.ListenAndServe())
	return nil
}
