package model


type Config struct{
	AppPort string   `json:"app_port"`
	DataBase DataBaseConfig `json:"db_config"`
}

type DataBaseConfig struct{
	DataSource string `json:"data_source"`
	DataBaseName string `json:"database"`
	AlbumTable string `json:"album_table"`
	PhotoTable string `json:"photo_table"`
}