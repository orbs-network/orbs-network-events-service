package events

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type storage struct {
	db *sql.DB
}

type Storage interface {
	StoreEvent(event *codec.Event) error
	GetEvents(contractName string, eventType string) ([]*codec.Event, error)
}

func NewStorage(dataSource string) (Storage, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}

	return &storage{
		db,
	}, nil
}

func (s *storage) StoreEvent(event *codec.Event) error {
	tableName := getTableName(event)

	if !s.checkIfTableExists(tableName) {
		if err := s.createTable(tableName, getTableMapping(event)); err != nil {
			return err
		}
	}

	return s.insertValues(tableName, event.Arguments)
}

func (s *storage) checkIfTableExists(tableName string) bool {
	row := s.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='" + tableName + "';")

	var count int
	row.Scan(&count)

	return count == 1
}

func (s *storage) createTable(tableName string, mapping string) error {
	query := "CREATE TABLE " + escape(tableName) + " (" + mapping + ")"
	println(query)
	_, err := s.db.Exec(query)
	return err
}

func (s *storage) insertValues(tableName string, values []interface{}) error {
	return errors.New("not implemented")
}

func (s *storage) GetEvents(contractName string, eventType string) ([]*codec.Event, error) {
	var events []*codec.Event

	query := "SELECT * FROM " + escape(getTableNameForContractAndEvent(contractName, eventType))
	println(query)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	columnsCount := len(columns)
	println(fmt.Sprintf("%+v", columns))

	for rows.Next() {
		arguments := make([]interface{}, columnsCount)
		if err := rows.Scan(&arguments); err != nil {
			return nil, err
		}

		println("RRR", arguments)
	}

	return events, nil
}

func getTableName(event *codec.Event) string {
	return getTableNameForContractAndEvent(event.ContractName, event.EventName)
}

func getTableNameForContractAndEvent(contractName string, eventName string) string {
	return strings.Join([]string{"events", contractName, eventName}, ".")
}

func getTableMapping(event *codec.Event) string {
	columns := []string{
		"timestamp uint64",
	}

	for index, arg := range event.Arguments {
		columns = append(columns, "argument"+strconv.FormatInt(int64(index), 10)+" "+getTableMappingType(arg))
	}

	return strings.Join(columns, ", ")
}

func getTableMappingType(arg interface{}) string {
	switch arg.(type) {
	case string:
		return "string"
	case uint32:
		return "uint32"
	default:
		return "blob"
	}
}

func escape(tableName string) string {
	return `"` + tableName + `"`
}
