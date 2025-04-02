package db

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON []byte

type Path struct {
	Sequences []PathSegment
}

type PathSegment struct {
	Seq     int
	Node    int64
	Edge    int64
	AggCost float64
}

func (p Path) ToTable() string {
	expr := "(VALUES "
	n := len(p.Sequences)

	for i, seq := range p.Sequences {
		expr += fmt.Sprintf("(%d,%d,%d,%f)", seq.Seq, seq.Node, seq.Edge, seq.AggCost)
		if i < n-1 {
			expr += ","
		}
	}

	return expr + ") AS path(seq,node,edge,agg_cost)"
}

func (p Path) ToArray() string {
	expr := "{"
	n := len(p.Sequences)

	for i, seq := range p.Sequences {
		if i > 0 && i < n-1 {
			expr += fmt.Sprintf("%d,", seq.Node)
		}
	}

	return expr + "}"
}

type GeoJSON struct {
	Type      string       `json:"type"`
	TotalCost float64      `json:"total_cost"`
	Features  []GeoFeature `json:"features"`
}

func (g GeoJSON) ToBytes(res http.ResponseWriter) []byte {
	bytes, err := json.Marshal(g)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	return bytes
}

var EmptyGeoJSON GeoJSON = GeoJSON{
	Type:     "FeatureCollection",
	Features: nil,
}

type GeoFeature struct {
	Type       string        `json:"type"`
	Geometry   GeoGeometry   `json:"geometry"`
	Properties GeoProperties `json:"properties"`
}

type GeoGeometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type GeoProperties struct {
	Gid        int     `json:"gid"`
	StreetName string  `json:"street_name"`
	RouteNum   string  `json:"route_num"`
	Fow        int8    `json:"fow"`
	AngleDiff  float64 `json:"angle_diff"`
	Distance   float64 `json:"distance"`
	Duration   float64 `json:"duration"`
}
