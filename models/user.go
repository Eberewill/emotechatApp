// models/user.go
package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`

	ConversationsSent     []Conversation `gorm:"foreignKey:SenderID"`   // Conversations where the user is the sender
	ConversationsReceived []Conversation `gorm:"foreignKey:ReceiverID"` // Conversations where the user is the receiver
}
