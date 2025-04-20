package types

import (
	"encoding/json"
	"time"
)

type GeoJSON struct {
	Type         string        `json:"type"`
	TotalCost    float64       `json:"total_cost"`
	Features     []GeoFeature  `json:"features"`
	QueryTime    time.Duration `json:"query_time_ns"`
	ResponseTime time.Duration `json:"response_time_ns"`
}

func (g *GeoJSON) ToBytes() []byte {
	bytes, _ := json.Marshal(g)
	return bytes
}

var EmptyGeoJSON GeoJSON = GeoJSON{
	Type:         "FeatureCollection",
	TotalCost:    -1,
	Features:     nil,
	QueryTime:    time.Duration(0),
	ResponseTime: time.Duration(0),
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
	StreetName string  `json:"street_name"`
	Fow        int8    `json:"fow"`
	AngleDiff  float64 `json:"angle_diff"`
	Distance   float64 `json:"distance"`
	Duration   float64 `json:"duration"`
}
