package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleRoute(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateRouteParams(params) {
		res.Write([]byte(`{"type" : "FeatureCollection", "features" : null}`))
		return
	}

	from, _ := strconv.Atoi(params.Get("from"))
	to, _ := strconv.Atoi(params.Get("to"))

	route := Db.Route(from, to)

	jsonBytes, err := json.Marshal(route)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = res.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}

func Generate() {
	Db.GenerateTNR(8)
}

func HandleGeocode(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateGeocodeParams(params) {
		res.Write([]byte("[]"))
		return
	}

	res.Write(Db.Geocode(params.Get("address")))
}
