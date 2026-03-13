package repository

import (
	"context"
	"fmt"

	"nunu/internal/model"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type MemoryRepository interface {
	Create(ctx context.Context, memory *model.Memory) error
	SearchByEmbedding(ctx context.Context, embedding []float32, limit int) ([]model.Memory, error)
}

type memoryRepository struct {
	db *gorm.DB
}

func NewMemoryRepository(db *gorm.DB) MemoryRepository {
	return &memoryRepository{db: db}
}

func (r *memoryRepository) Create(ctx context.Context, memory *model.Memory) error {
	return r.db.WithContext(ctx).Create(memory).Error
}

func (r *memoryRepository) SearchByEmbedding(ctx context.Context, embedding []float32, limit int) ([]model.Memory, error) {
	var memories []model.Memory

	vec := pgvector.NewVector(embedding)
	err := r.db.WithContext(ctx).
		Select("*, embedding <=> ? AS distance", vec).
		Order(fmt.Sprintf("embedding <=> '%s'", vec.String())).
		Limit(limit).
		Find(&memories).Error
	if err != nil {
		return nil, fmt.Errorf("memory search error: %w", err)
	}

	return memories, nil
}
