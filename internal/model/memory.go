package model

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type Memory struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	ConversationID string          `gorm:"index;not null" json:"conversation_id"`
	UserQuestion   string          `gorm:"type:text;not null" json:"user_question"`
	Summary        string          `gorm:"type:text" json:"summary"`
	SQLUsed        string          `gorm:"type:text" json:"sql_used"`
	ToolsUsed      string          `gorm:"type:text" json:"tools_used"` // JSON array
	Embedding      pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (Memory) TableName() string {
	return "memories"
}
