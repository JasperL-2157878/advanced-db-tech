package handlers

import (
	"strconv"

	"example.com/source/types"
)

func HandleGeocode(ctx Context) {
	if !ctx.validateGeocodeParams() {
		ctx.res.Write(types.JSON("[]"))
		return
	}

	address := ctx.Param("address")
	geocoding := ctx.db.Geocode(address)

	ctx.res.Write(geocoding)
}

func parseRouteRequest(ctx Context) (from int64, to int64) {
	if !ctx.validateRouteParams() {
		ctx.res.Write(types.EmptyGeoJSON.ToBytes())
		return
	}

	f, _ := strconv.Atoi(ctx.Param("from"))
	t, _ := strconv.Atoi(ctx.Param("to"))

	return int64(f), int64(t)
}

func HandleAlgDijkstra(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.Dijkstra(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleAlgAstar(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.Astar(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleAlgBdDijkstra(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.BdDijkstra(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleAlgBdAstar(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.BdAstar(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleOptNone(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.graphs.Base.BdDijkstra(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleOptTnr(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.graphs.Tnr.BdDijkstra(from, to))

	ctx.res.Write(route.ToBytes())
}

/*func HandleCh(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.Ch(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleTnrCh(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.db.TnrCh(from, to))

	ctx.res.Write(route.ToBytes())
} */
