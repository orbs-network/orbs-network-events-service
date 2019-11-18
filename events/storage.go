package events

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
)

type storage struct {
	db *sql.DB
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

func (s *storage) StoreEvent(blockHeight uint64, timestamp int64, event *codec.Event) error {
	tableName := getTableName(event)

	if !s.checkIfTableExists(tableName) {
		if err := s.createTable(tableName, getTableMapping(event)); err != nil {
			return err
		}
	}

	if columns, err := s.getTableColumns(tableName); err != nil {
		return err
	} else {
		values := []interface{}{
			blockHeight, timestamp,
		}
		values = append(values, event.Arguments...)
		return s.insertValues(tableName, columns, values...)
	}
}

func (s *storage) GetEvents(filterQuery *FilterQuery) (events []*StoredEvent, err error) {
	query := "SELECT * FROM " + escape(getTableNameForContractAndEvent(filterQuery.ContractName, filterQuery.EventNames[0]))
	println(query)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	columnsCount := len(columns)
	println("columns", fmt.Sprintf("%+v", columns))

	for rows.Next() {
		arguments := make([]interface{}, columnsCount)
		dest := make([]interface{}, columnsCount) // A temporary interface{} slice
		for i, _ := range arguments {
			dest[i] = &arguments[i]
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		println(fmt.Sprintf("args %+v", arguments))

		events = append(events, &StoredEvent{
			ContractName: filterQuery.ContractName,
			EventName:    filterQuery.EventNames[0],

			BlockHeight: uint64(arguments[0].(int64)),
			Timestamp:   arguments[1].(int64),
			Arguments:   arguments[2:],
		})
	}

	return events, nil
}

func (s *storage) GetBlockHeight() uint64 {
	panic("not implemented")
	row := s.db.QueryRow("SELECT count(*) FROM ;")

	var count uint64
	row.Scan(&count)

	return count
}

func (s *storage) IncBlockHeight() uint64 {
	panic("not implemented")
}
