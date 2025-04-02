package handlers

import (
	"net/http"
	"strconv"

	db "example.com/source/database"
)

func HandleGeocode(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	if !validateGeocodeParams(params) {
		res.Write(db.JSON("[]"))
		return
	}

	res.Write(Db.Geocode(params.Get("address")))
}

func parseRouteRequest(res http.ResponseWriter, req *http.Request) (from int64, to int64) {
	params := req.URL.Query()

	if !validateRouteParams(params) {
		res.Write(db.EmptyGeoJSON.ToBytes(res))
		return
	}

	f, _ := strconv.Atoi(params.Get("from"))
	t, _ := strconv.Atoi(params.Get("to"))

	return int64(f), int64(t)
}

func HandleDijkstra(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.Dijkstra(from, to))

	res.Write(route.ToBytes(res))
}

func HandleAstar(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.Astar(from, to))

	res.Write(route.ToBytes(res))
}

func HandleBdDijkstra(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.BdDijkstra(from, to))

	res.Write(route.ToBytes(res))
}

func HandleBdAstar(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.BdAstar(from, to))

	res.Write(route.ToBytes(res))
}

/*func HandleTnr(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.Tnr(from, to))

	res.Write(route.ToBytes(res))
}

func HandleCh(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.Ch(from, to))

	res.Write(route.ToBytes(res))
}

func HandleTnrCh(res http.ResponseWriter, req *http.Request) {
	from, to := parseRouteRequest(res, req)
	route := Db.Route(Db.TnrCh(from, to))

	res.Write(route.ToBytes(res))
} */
