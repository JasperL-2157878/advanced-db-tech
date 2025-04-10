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
	route := ctx.db.Route(ctx.graph.Base(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleOptTnr(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.graph.BaseTnr(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleOptCh(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.graph.Ch(from, to))

	ctx.res.Write(route.ToBytes())
}

func HandleOptChTnr(ctx Context) {
	from, to := parseRouteRequest(ctx)
	route := ctx.db.Route(ctx.graph.ChTnr(from, to))

	ctx.res.Write(route.ToBytes())
}
