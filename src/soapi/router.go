package soapi

import (
	"bikelib"
	"encoding/json"
	"maplib"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// chanStruct
type chanStruct struct {
	Error error
	Data  interface{}
	Type  string
}

// jsonResponse
func jsonResponse(w http.ResponseWriter, data interface{}) {
	resData, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Json marshal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(resData)
}

// Index router
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	router := make(map[string]string)

	router["/"] = "return this page"
	router["/address/:addr"] = "return gaode address object"
	router["/bike/:lng/:lat"] = "return near bike"

	data := ResIndex{
		Status: 1,
		Path:   router,
	}

	jsonResponse(w, data)
}

// Address router
func Address(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var resAddr ResAddr
	dataChan := make(chan chanStruct, 1)
	mapObj := maplib.MapAddr{
		Address: ps.ByName("addr"),
	}

	go func() {
		mapLoc, err := mapObj.GetGeoLoc()
		if err != nil {
			dataChan <- chanStruct{
				Error: err,
				Data:  nil,
			}
		} else {
			dataChan <- chanStruct{
				Error: nil,
				Data:  mapLoc,
			}
		}
	}()

	select {
	case chanRes := <-dataChan:
		if chanRes.Data == nil {
			close(dataChan)
			return
		}
		resData := chanRes.Data.(maplib.MapLocation)
		resAddr = ResAddr{
			Lat:      resData.Lat,
			Lng:      resData.Lng,
			Address:  resData.Address,
			Citycode: resData.CityCode,
			City:     resData.City,
		}

		log.WithFields(log.Fields{
			"Addr": resAddr.Address,
		}).Info("Find address")

		jsonResponse(w, resAddr)
	}

	close(dataChan)
}

// NearbyBike router
func NearbyBike(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lng, err := strconv.ParseFloat(ps.ByName("lng"), 64)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse float error")
		return
	}

	lat, err := strconv.ParseFloat(ps.ByName("lat"), 64)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse float error")
		return
	}

	dataChan := make(chan chanStruct, 2)
	go func() {
		bike := bikelib.Mobike{
			Lat: lat,
			Lng: lng,
		}
		data, err := bike.GetNearbyCar()
		if err != nil {
			dataChan <- chanStruct{
				Error: err,
				Data:  nil,
				Type:  "mobike",
			}
		} else {
			dataChan <- chanStruct{
				Error: nil,
				Data:  data,
				Type:  "mobike",
			}
		}
	}()

	go func() {
		bike := bikelib.Ofobike{
			Lat: lat,
			Lng: lng,
		}
		data, err := bike.GetNearbyCar()
		if err != nil {
			dataChan <- chanStruct{
				Error: err,
				Data:  nil,
				Type:  "ofo",
			}
		} else {
			dataChan <- chanStruct{
				Error: nil,
				Data:  data,
				Type:  "ofo",
			}
		}
	}()

	var resData ResBike
	for i := 0; i < cap(dataChan); i++ {
		select {
		case bikeChan := <-dataChan:
			generateBikeData(bikeChan, &resData)
		}
	}
	close(dataChan)

	jsonResponse(w, resData)
}

func generateBikeData(bikeChan chanStruct, data *ResBike) {
	if bikeChan.Error != nil {
		log.WithFields(log.Fields{
			"Error": bikeChan.Error,
		}).Error("Bike data contains error!")
		data.Status = 0
		return
	}

	if bikeChan.Type == "ofo" {
		bikeList := bikeChan.Data.([]bikelib.OfoCar)

		data.Ofo.Count = len(bikeList)
		log.WithFields(log.Fields{
			"Ofo": data.Ofo.Count,
		}).Info("Find ofo cars")

		data.Ofo.Data = make([]BikeEl, len(bikeList))
		for index, el := range bikeList {
			data.Ofo.Data[index].Lat = el.Lat
			data.Ofo.Data[index].Lng = el.Lng
			data.Ofo.Data[index].ID = el.Carno
			data.Ofo.Data[index].Type = "ofo"
		}
	}

	if bikeChan.Type == "mobike" {
		bikeList := bikeChan.Data.([]bikelib.MobikeCar)

		data.Mobike.Count = len(bikeList)
		log.WithFields(log.Fields{
			"Mobike": data.Mobike.Count,
		}).Info("Find mobike cars")

		data.Mobike.Data = make([]BikeEl, len(bikeList))
		for index, el := range bikeList {
			data.Mobike.Data[index].Lat = el.DistY
			data.Mobike.Data[index].Lng = el.DistX
			data.Mobike.Data[index].ID = el.DistID
		}
	}
	data.Status = 1
}

// RedirectHome func
func RedirectHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/web/", 302)
}
