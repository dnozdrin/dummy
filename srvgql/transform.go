package srvgql

import (
	"github.com/akhripko/dummy/models"
	gqlmodels "github.com/akhripko/dummy/srvgql/models"
)

func helloMessageToMessage(message *models.HelloMessage) (*gqlmodels.Message, error) {
	return &gqlmodels.Message{
		Data: message.Message,
	}, nil

}
