package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func HandleRoute(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateRouteParams(params) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	from, _ := strconv.Atoi(params.Get("from"))
	to, _ := strconv.Atoi(params.Get("to"))

	res.Write(Db.Route(from, to))
}

func HandleGeocode(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateGeocodeParams(params) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	street, number, postal, city := parseAddress(params.Get("address"))

	res.Write(Db.Geocode(street, number, postal, city))
}

func parseAddress(address string) (street string, number int, postal string, city string) {
	address = strings.ReplaceAll(address, ",", "")
	parts := strings.Split(address, " ")

	if len(parts) > 0 {
		street = parts[0]
	}
	if len(parts) > 1 {
		number, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		postal = parts[2]
	}
	if len(parts) > 3 {
		city = parts[3]
	}

	return
}
