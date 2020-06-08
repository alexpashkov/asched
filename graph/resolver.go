package graph

import "github.com/alexpashkov/asched/internal/amenities"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AmenitiesService *amenities.Service
}
