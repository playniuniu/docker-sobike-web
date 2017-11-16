package soapi 

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

// BikeEl struct
type BikeEl struct {
	Lng  float64 `json:"lng"`
	Lat  float64 `json:"lat"`
	Type string  `json:"type"`
	ID   string  `json:"id"`
}

// ResBike struct
type ResBike struct {
	Status int `json:"status"`
	Ofo    struct {
		Data  []BikeEl `json:"bike_list"`
		Count int      `json:"count"`
	} `json:"ofo"`
	Mobike struct {
		Data  []BikeEl `json:"bike_list"`
		Count int      `json:"count"`
	} `json:"mobike"`
}
