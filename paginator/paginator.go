package paginator

import (
	"context"
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * Created by Muhammad Muflih Kholidin
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

// MPaginator holds pagination options
type MPaginator struct {
	Collection *mongo.Collection
	Filter     bson.D
	Page       int64
	Size       int64
	Sort       map[string]int
}

// FindOptions creates and returns a configured *options.FindOptions
func FindOptions(page, size int64, sortField string, direction int) *options.FindOptions {
	// Ensure valid direction: 1 for ascending, -1 for descending
	if direction != 1 && direction != -1 {
		direction = 1 // Default to ascending if invalid
	}

	// Create sort document
	sort := bson.D{{Key: sortField, Value: direction}}

	// Calculate skip for pagination
	skip := (page - 1) * size

	return options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(size)
}

// GetPaginator calculates pagination details and returns paginated results
func (p *MPaginator) GetPaginator(ctx context.Context, results interface{}) (*PaginatorResponse, error) {
	// Define the pagination options
	field := ""
	direction := -1
	for v, k := range p.Sort {
		field = v
		direction = k
	}
	findOptions := FindOptions(p.Page, p.Size, field, direction)

	// Get total count of documents
	totalCount, err := p.Collection.CountDocuments(ctx, p.Filter)
	if err != nil {
		return nil, fmt.Errorf("error counting documents: %w", err)
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalCount) / float64(p.Size)))

	// Perform the query
	cursor, err := p.Collection.Find(ctx, p.Filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	// Determine previous and next pages
	var prevPage, nextPage *int64
	if p.Page > 1 {
		prev := p.Page - 1
		prevPage = &prev
	}
	if p.Page < totalPages {
		next := p.Page + 1
		nextPage = &next
	}

	return &PaginatorResponse{
		Data: results,
		Paginate: &PaginatorData{
			Count:      totalCount,
			Page:       p.Page,
			Size:       p.Size,
			TotalPages: totalPages,
			NextPage:   nextPage,
			PrevPage:   prevPage,
		},
	}, nil
}

// PaginatorResponse holds paginated data and metadata
type PaginatorResponse struct {
	Data     interface{}    `json:"data"`
	Paginate *PaginatorData `json:"paginate"`
}

// PaginatorData holds pagination metadata
type PaginatorData struct {
	Count      int64  `json:"total"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
	TotalPages int64  `json:"total_pages"`
	NextPage   *int64 `json:"next_page"`
	PrevPage   *int64 `json:"prev_page"`
}
