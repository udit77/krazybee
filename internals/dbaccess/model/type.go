package model

type AlbumData struct{
	ID int `json:"id"`
	UserID int `json:"userId"`
	Title  string `json:"title"`
}


type PhotoData struct{
	ID int `json:"id"`
	AlbumId int `json:"albumId"`
	Title  string `json:"title"`
	Url  string `json:"url"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
}