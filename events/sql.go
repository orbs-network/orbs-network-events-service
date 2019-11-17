package events

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"strconv"
	"strings"
)

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

func (s *storage) insertValues(tableName string, columns []string, values ...interface{}) error {
	var fillIns []string
	for i := 0; i < len(columns); i++ {
		fillIns = append(fillIns, "?")
	}

	query := "INSERT INTO " + escape(tableName) + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(fillIns, ", ") + ")"
	println(query)
	statement, _ := s.db.Prepare(query)
	println(fmt.Sprintf("%+v", values))
	_, err := statement.Exec(values...)
	return err
}

func getTableName(event *codec.Event) string {
	return getTableNameForContractAndEvent(event.ContractName, event.EventName)
}

func getTableNameForContractAndEvent(contractName string, eventName string) string {
	return strings.Join([]string{"events", contractName, eventName}, ".")
}

func getTableMapping(event *codec.Event) string {
	columns := []string{
		"blockHeight uint64",
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

func (s *storage) getTableColumns(tableName string) ([]string, error) {
	query := "SELECT * FROM " + escape(tableName) + " LIMIT 0"
	println(query)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows.Columns()
}
