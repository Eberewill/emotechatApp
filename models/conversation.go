// models/conversation.go
package models

import "time"

type Conversation struct {
	ID         uint      `gorm:"primaryKey"`
	ReceiverID uint      // ID of the logged-in user (receiver)
	SenderID   uint      // ID of the conversation partner (sender)
	Content    string    // Last message content
	Timestamp  time.Time // Last message timestamp

	Receiver User `gorm:"foreignKey:ReceiverID"` // Receiver (logged-in user)
	Sender   User `gorm:"foreignKey:SenderID"`   // Sender (conversation partner)
}
