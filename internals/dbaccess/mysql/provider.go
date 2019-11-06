package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/dbaccess"
	"github.com/krazybee/internals/dbaccess/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Provider struct{
	config config.Provider
}

func NewDBProvider(provider config.Provider) (dbaccess.Provider, error) {
	dbConfig := provider.GetConfig().DataBase
	dataSource := dbConfig.DataSource + "/"
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS "+dbConfig.DataBaseName)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("USE "+dbConfig.DataBaseName)
	if err != nil {
		return nil, err
	}

	_,err = db.Exec("DROP TABLE IF EXISTS  " + dbConfig.PhotoTable)
	if err != nil {
		return nil, err
	}

	_,err = db.Exec("DROP TABLE IF EXISTS  " + dbConfig.AlbumTable)
	if err != nil {
		return nil, err
	}

	_,err = db.Exec("CREATE TABLE " + dbConfig.AlbumTable + " ( id integer NOT NULL, userId integer, title text, PRIMARY KEY (id) ) ")
	if err != nil {
		return nil, err
	}
	_,err = db.Exec("CREATE TABLE " + dbConfig.PhotoTable + " ( id integer, albumId integer NOT NULL, title text, url text, thumbnailUrl text, FOREIGN KEY (albumId) REFERENCES "+ dbConfig.AlbumTable + "(id))")
	if err != nil {
		return nil, err
	}

	err = initializeData(db,dbConfig.AlbumTable,dbConfig.PhotoTable)
	if err != nil {
		return nil,err
	}
	return &Provider{config:provider}, nil
}


func (i *Provider) GetDBConn() *sql.DB {
	dbConfig := i.config.GetConfig().DataBase
	db, err := sql.Open("mysql", dbConfig.DataSource + "/"+dbConfig.DataBaseName)
	if err != nil {
		log.Fatal("[MySql] error getting db connection")
	}
	return db
}


func getAlbums() ([]model.AlbumData, error){
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	var albums []model.AlbumData
	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/albums", nil)
	if err != nil {
		return albums, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return albums, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return albums, err
	}

	respContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return albums, err
	}
	err = json.Unmarshal(respContent,&albums)
	return albums, nil
}

func getPhotoDetails(album model.AlbumData) ([]model.PhotoData, error){
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	var photos []model.PhotoData
	req, err := http.NewRequest("GET", fmt.Sprintf("%v%v","https://jsonplaceholder.typicode.com/photos?albumId=",album.ID), nil)
	if err != nil {
		return photos, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return photos, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return photos, err
	}

	respContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return photos, err
	}
	err = json.Unmarshal(respContent,&photos)
	return photos, nil
}

func initializeData(db *sql.DB, albumTable string, photoTable string) error{
	log.Println("[InitializeData] initializing data from api")
	albums, err := getAlbums()
	if err != nil{
		log.Println("[MySQL] error fetching albums data", err)
		return err
	}
	for i:=0; i<len(albums); i++ {
		insertQuery, err := db.Prepare("INSERT INTO "+ albumTable+"(id, userId, title) VALUES(?,?,?)")
		if err != nil {
			return err
		}
		_, err = insertQuery.Exec(albums[i].ID, albums[i].UserID,albums[i].Title)
		if err != nil{
			return err
		}
		photos, err := getPhotoDetails(albums[i])
		if err != nil {
			return err
		}
		for i:=0;i<len(photos);i++{
			insertQuery, err := db.Prepare("INSERT INTO "+ photoTable+"(id, albumId, title, url, thumbnailUrl) VALUES(?,?,?,?,?)")
			if err != nil {
				return err
			}
			_, err = insertQuery.Exec(photos[i].ID, photos[i].AlbumId, photos[i].Title, photos[i].Url, photos[i].ThumbnailUrl)
			if err != nil{
				return err
			}
		}
	}
	log.Println("[InitializeData] data initialization successful")
	return nil
}
