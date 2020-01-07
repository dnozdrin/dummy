package srvgql

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/akhripko/dummy/models"
	gqlmodels "github.com/akhripko/dummy/srvgql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	service Service
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

var internalSrvErr models.ErrInternalSrvErr

func (r *queryResolver) Hello(_ context.Context, name string) (*gqlmodels.Message, error) {
	message, err := r.service.Hello(name)
	if err != nil {
		switch err.(type) {
		case models.ErrNotValidRequest:
			return &gqlmodels.Message{
				Error: &gqlmodels.Error{
					Code:    1,
					Message: err.Error(),
				},
			}, nil
		default:
			log.Error("gqlsrv: failed to build service token:", err)
			return nil, internalSrvErr
		}
	}
	return helloMessageToMessage(message)
}
