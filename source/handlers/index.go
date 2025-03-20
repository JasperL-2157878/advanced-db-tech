package handlers

import (
	"net/http"
	"strconv"
)

func HandleRoute(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateRouteParams(params) {
		println("valid")
		res.Write([]byte(`{"type" : "FeatureCollection", "features" : null}`))
		return
	}

	from, _ := strconv.Atoi(params.Get("from"))
	to, _ := strconv.Atoi(params.Get("to"))

	res.Write(Db.Route(from, to))
}

func HandleGeocode(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateGeocodeParams(params) {
		println("invalid")
		res.Write([]byte("[]"))
		return
	}

	res.Write(Db.Geocode(params.Get("address")))
}
