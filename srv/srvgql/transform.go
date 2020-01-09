package srvgql

import (
	"github.com/akhripko/dummy/models"
	models2 "github.com/akhripko/dummy/srv/srvgql/models"
)

func helloMessageToMessage(message *models.HelloMessage) (*models2.Message, error) {
	return &models2.Message{
		Data: message.Message,
	}, nil

}
