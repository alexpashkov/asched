package model

import "github.com/99designs/gqlgen/graphql"

type NewAmenity struct {
	Name        string          `json:"name"`
	Type        []string        `json:"type"`
	Lat         float64         `json:"lat"`
	Lon         float64         `json:"lon"`
	Photo       *graphql.Upload `json:"photo"`
	Description *string         `json:"description"`
}

type Amenity struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        []string `json:"type"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	Description *string  `json:"description"`
}
