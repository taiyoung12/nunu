package repository

import (
	"context"

	"nunu/internal/model"

	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(ctx context.Context, conv *model.Conversation) error
	Update(ctx context.Context, conv *model.Conversation) error
	GetByThreadTS(ctx context.Context, channelID, threadTS string) (*model.Conversation, error)
}

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(ctx context.Context, conv *model.Conversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

func (r *conversationRepository) Update(ctx context.Context, conv *model.Conversation) error {
	return r.db.WithContext(ctx).Save(conv).Error
}

func (r *conversationRepository) GetByThreadTS(ctx context.Context, channelID, threadTS string) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND thread_ts = ?", channelID, threadTS).
		First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}
