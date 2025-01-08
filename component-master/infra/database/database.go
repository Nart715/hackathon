package database

import (
	"component-master/config"
	"fmt"
	"log/slog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSource struct {
	db *gorm.DB
}

type Transaction struct {
	tx *gorm.DB
}

func NewDataSource(conf *config.DataSourceConfig) (*DataSource, error) {
	dns := conf.BuildDnsMYSQL()
	slog.Info(fmt.Sprintf("connecting to database %s, %s", "url", dns))
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		IgnoreRelationshipsWhenMigrating: true,
		QueryFields:                      true,
		PrepareStmt:                      true,
		AllowGlobalUpdate:                true,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database, %v", err))
		return nil, err
	}
	instance, err := db.DB()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database, %v", err))
		return nil, err
	}
	instance.SetMaxIdleConns(conf.MaxIdleConnections)
	instance.SetMaxOpenConns(conf.MaxOpenConnections)
	instance.SetConnMaxIdleTime(conf.MaxConnIdleTime)
	return &DataSource{
		db: db,
	}, nil
}

func (db *DataSource) GetDB() *gorm.DB {
	return db.db
}

func (db *DataSource) Close() {
	slog.Info("closing database connection")
	instance, err := db.db.DB()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to close database connection, %v", err))
		panic(err)
	}
	instance.Close()
}
