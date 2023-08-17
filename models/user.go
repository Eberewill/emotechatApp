// models/user.go
package models

import "time"

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`

	ResetToken            string         `json:"-"`
	ResetTokenExpiry      time.Time      `json:"-"`
	ConversationsSent     []Conversation `gorm:"foreignKey:SenderID"`   // Conversations where the user is the sender
	ConversationsReceived []Conversation `gorm:"foreignKey:ReceiverID"` // Conversations where the user is the receiver
}
