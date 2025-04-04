package types

import (
	"encoding/json"
)

type GeoJSON struct {
	Type      string       `json:"type"`
	TotalCost float64      `json:"total_cost"`
	Features  []GeoFeature `json:"features"`
}

func (g GeoJSON) ToBytes() []byte {
	bytes, _ := json.Marshal(g)
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
