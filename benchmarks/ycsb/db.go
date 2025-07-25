package ycsb

import (
	"context"
	"fmt"
)

type DBConfig struct {
	Params map[string]string
}

// DBCreator creates a database layer.
type DBCreator interface {
	Create() (DB, error)
}

// DB is the layer to access the database to be benchmarked.
type DB interface {
	// Close closes the database layer.
	Close() error

	// InitThread initializes the state associated to the goroutine worker.
	// The Returned context will be passed to the following usage.
	InitThread(ctx context.Context, threadID int, threadCount int) context.Context

	// CleanupThread cleans up the state when the worker finished.
	CleanupThread(ctx context.Context)

	// Read reads a record from the database and returns a map of each field/value pair.
	// table: The name of the table.
	// key: The record key of the record to read.
	// fields: The list of fields to read, nil|empty for reading all.
	Read(ctx context.Context, table string, key string) (string, error)

	// Update updates a record in the database. Any field/value pairs will be written into the
	// database or overwritten the existing values with the same field name.
	// table: The name of the table.
	// key: The record key of the record to update.
	// values: A map of field/value pairs to update in the record.
	Update(ctx context.Context, table string, key string, value string) error

	// Insert inserts a record in the database. Any field/value pairs will be written into the
	// database.
	// table: The name of the table.
	// key: The record key of the record to insert.
	// values: A map of field/value pairs to insert in the record.
	Insert(ctx context.Context, table string, key string, value string) error

	// Delete deletes a record from the database.
	// table: The name of the table.
	// key: The record key of the record to delete.
	Delete(ctx context.Context, table string, key string) error
}
type TransactionDB interface {
	DB

	// NewTransaction creates a new datastore based on itself with a clean state.
	NewTransaction() TransactionDB

	// Start starts a new transaction.
	Start() error

	// Commit commits the current transaction.
	Commit() error

	// Abort aborts the current transaction.
	Abort() error
}

type BatchDB interface {
	// BatchInsert inserts batch records in the database.
	// table: The name of the table.
	// keys: The keys of batch records.
	// values: The values of batch records.
	BatchInsert(ctx context.Context, table string, keys []string, values []map[string][]byte) error

	// BatchRead reads records from the database.
	// table: The name of the table.
	// keys: The keys of records to read.
	// fields: The list of fields to read, nil|empty for reading all.
	BatchRead(
		ctx context.Context,
		table string,
		keys []string,
		fields []string,
	) ([]map[string][]byte, error)

	// BatchUpdate updates records in the database.
	// table: The name of table.
	// keys: The keys of records to update.
	// values: The values of records to update.
	BatchUpdate(ctx context.Context, table string, keys []string, values []map[string][]byte) error

	// BatchDelete deletes records from the database.
	// table: The name of the table.
	// keys: The keys of the records to delete.
	BatchDelete(ctx context.Context, table string, keys []string) error
}

var dbCreators = map[string]DBCreator{}

// RegisterDBCreator registers a creator for the database
func RegisterDBCreator(name string, creator DBCreator) {
	_, ok := dbCreators[name]
	if ok {
		panic(fmt.Sprintf("duplicate register database %s", name))
	}

	dbCreators[name] = creator
}

// GetDBCreator gets the DBCreator for the database
func GetDBCreator(name string) DBCreator {
	return dbCreators[name]
}
