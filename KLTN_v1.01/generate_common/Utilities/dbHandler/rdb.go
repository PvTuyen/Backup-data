package dbHandler

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DriverNameMaster = "postgres"
	DriverNameTarget = "postgres"
)
type sql interface {
	QuerySql(sql string) error
	GetListTable() error
	GetListColumns(nameTable string) error
	GetListKeys(nameTable string) error
}

type sqlDb struct {
	sqlx *sqlx.DB
}
type SchemaDb struct {
	ListNameTable []string
	ListColumns   map[string][]string
	ListKeys      map[string][]string
	ListType map[string][]string
}


var (
	sqlDbInstance sqlDb
	schemaDbInstance = SchemaDb{
		ListNameTable: []string{},
		ListColumns:   make(map[string][]string),
		ListKeys:      make(map[string][]string),
		ListType: make(map[string][]string),
	}
)

//type SqlQuery interface {
//	Generate(query string) error
//}
//
//type DbSqlx struct {
//	sqlx *sqlx.DB
//}
//
//func (rdbRepository DbSqlx) Generate(query string) error {
//	_, err := rdbRepository.sqlx.Exec(query)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//


func (con *sqlDb) QuerySql(sql string) error  {
	_, err := con.sqlx.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func (con *sqlDb) GetListTable()  error {
	if con.sqlx == nil {
	fmt.Println("== nil")
	}
	rows, err := con.sqlx.Query("SELECT table_name FROM INFORMATION_SCHEMA.TABLES where table_schema = 'public';")
	if err != nil {
		fmt.Println("dbHandler GetListTable connect failed", err.Error())
		return err
	}
	defer rows.Close()
	names := []string{}
	checkScan := true
	for rows.Next() {
		var name string
		if errS := rows.Scan(&name); errS != nil {
			fmt.Println("dbHandler GetListTable scan failed")
			checkScan = false
			break
		}
		names = append(names, name)
	}
	if checkScan {
		schemaDbInstance.ListNameTable = names
		//fmt.Println(schemaDbInstance.ListNameTable)
		return nil
	} else {
		return err
	}
}
func (con *sqlDb) GetListColumns(nameTable string)  error  {
	//fmt.Println("name table",nameTable)
	if con.sqlx == nil {
		fmt.Println("== nil")
	}
	sqlQ := fmt.Sprintf("SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%v';", nameTable)
	rows, err := con.sqlx.Query(sqlQ)
	if err != nil {
		fmt.Println("dbHandler GetListTable connect failed")
		return err
	}
	defer rows.Close()
	names := []string{}
	types := []string{}
	checkScan := true
	for rows.Next() {
		var name string
		var typeC string
		if errS := rows.Scan(&name, &typeC); errS != nil {
			fmt.Println("dbHandler GetListTable scan failed")
			checkScan = false
			break
		}
		names = append(names, name)
		types = append(types, typeC)

	}
	//fmt.Println(names)
	if checkScan {
		schemaDbInstance.ListColumns[nameTable] = names
		schemaDbInstance.ListType[nameTable] = types
		//fmt.Println(schemaDbInstance.ListColumns[nameTable])
		return nil
	} else {
		return err
	}
}
func (con *sqlDb) GetListKeys(nameTable string) error  {
	//fmt.Println("name table",nameTable)
	if con.sqlx == nil {
		fmt.Println("== nil")
	}
	sqlQ := fmt.Sprintf("SELECT column_name FROM INFORMATION_SCHEMA.key_column_usage WHERE TABLE_NAME = '%v';", nameTable)
	rows, err := con.sqlx.Query(sqlQ)
	if err != nil {
		fmt.Println("dbHandler GetListKeys connect failed")
		return err
	}
	defer rows.Close()
	names := []string{}
	checkScan := true
	for rows.Next() {
		var name string
		if errS := rows.Scan(&name); errS != nil {
			fmt.Println("dbHandler GetListKeys scan failed")
			checkScan = false
			break
		}
		names = append(names, name)
	}
	//fmt.Println(names)
	if checkScan {
		schemaDbInstance.ListKeys[nameTable] = names
		//fmt.Println(schemaDbInstance.ListColumns[nameTable])
		return nil
	} else {
		return err
	}
}

func InitDB(driverName string, info string) error {
	con , err := sqlx.Open(driverName, info)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("dbHandler InitDb failed")

		return nil
	}
	sqlDbInstance.sqlx = con
	return err
}
func GetSchemaInstance() *SchemaDb {
	return &schemaDbInstance
}
func ConnectDb(driverName string, info string) (*sqlx.DB, error)  {
	err0 := errors.New("@connect database failed")
	con , err := sqlx.Open(driverName, info)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("dbHandler InitDb failed")
		return con, err0
	}
	return con, nil
}