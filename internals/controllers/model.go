package controllers


type SearchAlbumResponse struct{
	UserID int `json:"userId"`
	ID int `json:"id"`
	Title string `json:"title"`
}


type SearchPhotoResponse struct{
	ID int `json:"id"`
	AlbumID int `json:"albumId"`
	Title string `json:"title"`
	Url string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}
