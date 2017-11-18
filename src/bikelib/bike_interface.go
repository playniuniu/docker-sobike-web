package bikelib

// BikeData is base bike struct
type BikeData struct {
	Lng      float64 `json:"lng"`
	Lat      float64 `json:"lat"`
	CarNo    string  `json:"car_no"`
	CarType  string  `json:"car_type"`
	Distance int     `json:"car_distance"`
}

// BikeInterface is common bike interface
type BikeInterface interface {
	GetNearbyCar() ([]BikeData, error)
}
