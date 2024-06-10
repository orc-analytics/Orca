package datalayer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dlyr "github.com/predixus/pdb_framework/internal/datalayer"
)

type DummySQLDb struct {
	mock string
}

type MockDB struct {
	*dlyr.Db
	DB *DummySQLDb
}

func (d *MockDB) Connect() error {
	db := DummySQLDb{mock: "Mock"}
	d.DB = &db
	return nil
}

func TestCloseResetsDatabase(t *testing.T) {
	db := MockDB{}
	db.Connect()
	assert.NotNil(t, db.DB, "Database is succesfully mocked")
	// TODO: Get this test working
	// db.Close()
	// assert.Nil(t, db.DB, "Database succesfully cleared on `Close`")
}
