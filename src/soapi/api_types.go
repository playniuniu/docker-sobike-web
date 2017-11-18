package soapi

import "bikelib"

// ResIndex struct
type ResIndex struct {
	Status int               `json:"status"`
	Path   map[string]string `json:"path"`
}

// ResAddr struct
type ResAddr struct {
	City     string  `json:"city"`
	Citycode string  `json:"citycode"`
	Address  string  `json:"address"`
	Lng      float64 `json:"lng"`
	Lat      float64 `json:"lat"`
}

// ResBike struct
type ResBike struct {
	Status int `json:"status"`
	Ofo    struct {
		Data  []bikelib.BikeData `json:"bike_list"`
		Count int                `json:"count"`
	} `json:"ofo"`
	Mobike struct {
		Data  []bikelib.BikeData `json:"bike_list"`
		Count int                `json:"count"`
	} `json:"mobike"`
}
