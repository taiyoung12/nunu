package repository

import (
	"context"
	"fmt"

	"nunu/internal/model"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type KnowledgeRepository interface {
	Create(ctx context.Context, knowledge *model.Knowledge) error
	SearchByEmbedding(ctx context.Context, embedding []float32, category string, limit int) ([]model.Knowledge, error)
	SearchByKeyword(ctx context.Context, keyword string, category string, limit int) ([]model.Knowledge, error)
}

type knowledgeRepository struct {
	db *gorm.DB
}

func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository {
	return &knowledgeRepository{db: db}
}

func (r *knowledgeRepository) Create(ctx context.Context, knowledge *model.Knowledge) error {
	return r.db.WithContext(ctx).Create(knowledge).Error
}

func (r *knowledgeRepository) SearchByEmbedding(ctx context.Context, embedding []float32, category string, limit int) ([]model.Knowledge, error) {
	var knowledges []model.Knowledge

	vec := pgvector.NewVector(embedding)
	query := r.db.WithContext(ctx)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.
		Order(fmt.Sprintf("embedding <=> '%s'", vec.String())).
		Limit(limit).
		Find(&knowledges).Error
	if err != nil {
		return nil, fmt.Errorf("knowledge embedding search error: %w", err)
	}

	return knowledges, nil
}

func (r *knowledgeRepository) SearchByKeyword(ctx context.Context, keyword string, category string, limit int) ([]model.Knowledge, error) {
	var knowledges []model.Knowledge

	query := r.db.WithContext(ctx)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query = query.Where("title ILIKE ? OR content ILIKE ? OR tags ILIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	err := query.Limit(limit).Find(&knowledges).Error
	if err != nil {
		return nil, fmt.Errorf("knowledge keyword search error: %w", err)
	}

	return knowledges, nil
}
