package domain

import (
	infra "orcaness.com/api/app/infra"
)

type RepositoryBase struct{}

// Publish events
func (this *RepositoryBase) PublishEvents(events []EventBase) error {
	for _, evt := range events {
		infra.Db().Create(evt)
	}
	return nil

}
