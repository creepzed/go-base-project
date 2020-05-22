package cockroach_connection

import (
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
)

type Migrate struct {
	connection CockroachConnection
}

func NewMigrate(connection CockroachConnection) *Migrate {
	return &Migrate{connection: connection}
}

func (m *Migrate) AutoMigrateAll(tables ...interface{}) {
	db := m.connection.GetConnection()
	db = db.AutoMigrate(tables...)
	if db.Error != nil {
		log.WithError(db.Error).Fatal(db.Error.Error())
	}
}
