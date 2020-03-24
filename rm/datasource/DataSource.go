package datasource

import (
	"database/sql"
)

type DataSource struct {
	resourceGroupId string
	resourceId      string
	db              *sql.DB
}

func DataSourceFactory(resourceGroupId, resourceId string, db *sql.DB) *DataSource {
	return &DataSource{resourceGroupId, resourceId, db}
}

func (d *DataSource) GetResourceGroupId() string {
	return d.resourceGroupId
}

func (d *DataSource) GetResourceId() string {
	return d.resourceId
}

func (d *DataSource) Begin() (LocalTransaction, error) {
	//todo
	return nil, nil
}

func (d *DataSource) getSQLDB() *sql.DB {
	return d.db
}
