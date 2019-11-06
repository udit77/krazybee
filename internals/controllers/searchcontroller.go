package controllers

import (
	"encoding/json"
	"errors"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/dbaccess"
	"log"
	"net/http"
	"strconv"
)

type SearchController struct {
	configProvider config.Provider
	dbAccessor   dbaccess.Provider
}

func NewController(confProvider config.Provider, dbProvider dbaccess.Provider) (*SearchController, error){
	if confProvider == nil{
		return &SearchController{}, errors.New("[Controllers][NewController] config provider is nil")
	}
	if dbProvider == nil{
		return &SearchController{}, errors.New("[Controllers][NewController] db provider is nil")
	}
	return &SearchController{
		configProvider:confProvider,
		dbAccessor:dbProvider,
	}, nil
}

func(s *SearchController) Search(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	queryType, found := queryParams["type"]
	typeExists := found && len(queryType[0]) >= 1

	if typeExists {
		if queryType[0] == "album" {
			id, found := queryParams["id"]
			idExists := found && len(id[0]) >= 1
			if idExists {
				albumId , err := strconv.Atoi(id[0])
				if err != nil{
					log.Println("[Controllers][SearchControllers] error converting id to integer",err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				albumTable := s.configProvider.GetConfig().DataBase.AlbumTable
				db := s.dbAccessor.GetDBConn()
				defer db.Close()
				selDB, err := db.Query("SELECT id, userId, title FROM "+ albumTable + " WHERE id=?", albumId)
				if err != nil {
					log.Println("[Controllers][SearchControllers] error querying data from database",err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				album := SearchAlbumResponse{}
				for selDB.Next() {
					err = selDB.Scan(&album.ID, &album.UserID, &album.Title)
					if err != nil {
						log.Println("[Controllers][SearchControllers] error scanning data from database",err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				responseData, err := json.Marshal(album)
				if err != nil{
					log.Println("[Controllers][SearchControllers] error marshalling album response")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Println("[Controllers][SearchController] api response",string(responseData))
				w.Write(responseData)
				return
			}else {
				log.Println("[Controllers][SearchControllers] invalid request album id not found")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}else if queryType[0] == "photo"{
			id, found := queryParams["id"]
			idExists := found && len(id[0]) >= 1
			album, found := queryParams["album"]
			albumExists := found && len(album[0]) >= 1

			if idExists && albumExists {
				photoId, err := strconv.Atoi(id[0])
				if err != nil {
					log.Println("[Controllers][SearchControllers] error converting id to integer", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				albumId, err := strconv.Atoi(album[0])
				if err != nil {
					log.Println("[Controllers][SearchControllers] error converting album id to integer", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				photoTable := s.configProvider.GetConfig().DataBase.PhotoTable
				db := s.dbAccessor.GetDBConn()
				defer db.Close()
				selDB, err := db.Query("SELECT id, albumId, title, url, thumbnailUrl FROM "+photoTable+ " WHERE id=? AND albumId=?", photoId,albumId)
				if err != nil {
					log.Println("[Controllers][SearchControllers] error querying data from database", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				var response []SearchPhotoResponse
				for selDB.Next() {
					photo := SearchPhotoResponse{}
					err = selDB.Scan(&photo.ID, &photo.AlbumID, &photo.Title, &photo.Url, &photo.ThumbnailUrl)
					if err != nil {
						log.Println("[Controllers][SearchControllers] error scanning data from database", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					response = append(response,photo)
				}
				responseData, err := json.Marshal(response)
				if err != nil {
					log.Println("[Controllers][SearchControllers] error marshalling photo response")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Println("[Controllers][SearchController] api response",string(responseData))
				w.Write(responseData)
				return
			}else {
				log.Println("[Controllers][SearchControllers] invalid photo id or album id not found")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}else{
			log.Println("[Controllers][SearchControllers] invalid request type ", queryType[0])
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}else{
		log.Println("[Controllers][SearchControllers] request type not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}