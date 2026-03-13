package model

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type Knowledge struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Category  string          `gorm:"index;not null" json:"category"` // "schema" | "example_sql" | "glossary"
	Title     string          `gorm:"not null" json:"title"`
	Content   string          `gorm:"type:text;not null" json:"content"`
	Tags      string          `gorm:"type:text" json:"tags"` // JSON array
	Embedding pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	CreatedAt time.Time       `json:"created_at"`
}

func (Knowledge) TableName() string {
	return "knowledges"
}
