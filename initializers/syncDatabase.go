package initializers

import (
	"github.com/Eberewill/emotechat/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Conversation{})
}
