package broker

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Broker struct {
	postgresDBSetup pgSetup
	postgresDB      *gorm.DB
}

func NewBroker() Broker {
	b := Broker{}
	return b
}

type pgSetup struct {
	username string
	password string
	dbName   string
	dbHost   string
}

func (b *Broker) InitialiseBroker() error {
	err := b.SetPostgres()
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) GetPostgres() *gorm.DB {
	return b.postgresDB
}

func (b *Broker) SetPostgres() error {
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		b.postgresDBSetup.dbHost,
		b.postgresDBSetup.username,
		b.postgresDBSetup.dbName,
		b.postgresDBSetup.password,
	)

	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}

	b.postgresDB = conn

	return nil
}

func (b *Broker) SetPostgresConfig(username, password, dbName, dbHost string) {
	pgs := pgSetup{
		username: username,
		password: password,
		dbName:   dbName,
		dbHost:   dbHost,
	}

	b.postgresDBSetup = pgs
}
