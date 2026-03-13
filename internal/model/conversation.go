package model

import "time"

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChannelID string    `gorm:"index;not null" json:"channel_id"`
	UserID    string    `gorm:"index;not null" json:"user_id"`
	ThreadTS  string    `gorm:"index" json:"thread_ts"`
	Question  string    `gorm:"type:text;not null" json:"question"`
	Answer    string    `gorm:"type:text" json:"answer"`
	SQLUsed   string    `gorm:"type:text" json:"sql_used"`
	ToolsUsed string    `gorm:"type:text" json:"tools_used"` // JSON array
	Success   bool      `gorm:"default:false" json:"success"`
	Feedback  string    `gorm:"type:varchar(20)" json:"feedback"` // "thumbsup" | "thumbsdown" | ""
	Duration  int64     `json:"duration_ms"`
	CreatedAt time.Time `json:"created_at"`
}

func (Conversation) TableName() string {
	return "conversations"
}
