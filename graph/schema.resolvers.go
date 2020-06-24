package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/alexpashkov/asched/graph/generated"
	"github.com/alexpashkov/asched/graph/model"
)

func (r *amenityResolver) Photos(_ context.Context, obj *model.Amenity) ([]string, error) {
	return r.AmenitiesService.GetPhotoIDs(obj.ID)
}

func (r *mutationResolver) AddAmenity(ctx context.Context, input model.NewAmenity) (string, error) {
	return r.AmenitiesService.AddAmenity(ctx, input)
}

func (r *queryResolver) Amenities(ctx context.Context, lat float64, lon float64, typeArg *string) ([]*model.Amenity, error) {
	return r.AmenitiesService.SearchAmenities(ctx, lat, lon, 100, typeArg)
}

// Amenity returns generated.AmenityResolver implementation.
func (r *Resolver) Amenity() generated.AmenityResolver { return &amenityResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type amenityResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
