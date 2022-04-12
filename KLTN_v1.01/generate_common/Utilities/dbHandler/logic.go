package dbHandler

import (
	"fmt"
	"generate_common/Utilities/tpl"
	"github.com/jmoiron/sqlx"
	"strings"
)

func Test() error {
	err := sqlDbInstance.GetListTable()
	if err != nil {
		return err
	}
	schemaD := schemaDbInstance.ListNameTable
	for _, v := range schemaD {
		err1 := sqlDbInstance.GetListColumns(v)
		if err1 != nil {
			return err1
		}
		err2 := sqlDbInstance.GetListKeys(v)
		if err2 != nil {
			return err2
		}
	}
	fmt.Println(schemaDbInstance.ListColumns)
	fmt.Println(schemaDbInstance.ListKeys)
	return nil
}

func ProcessSchema(schemaRoot *SchemaDb, schemaTarget *SchemaDb, dbRoot *sqlx.DB, dbTarget *sqlx.DB) {
	for _, v := range schemaRoot.ListNameTable{
		ProcessDataInTable(schemaRoot,schemaTarget,dbRoot,dbTarget,v)
		//break
	}
}
func UpdaterSql(schemaRoot *SchemaDb, schemaTarget *SchemaDb, dbRoot *sqlx.DB, dbTarget *sqlx.DB) {
	tableSnapshot := "snapshot"
	sqlx := "select * from "+ tableSnapshot + " order by time desc limit 1;"
	rows, err := dbRoot.Queryx(sqlx)
	if err != nil {
		fmt.Println("@ query failed")
		fmt.Println(sqlx)
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		values,err1 := rows.SliceScan()
		if values == nil || len(values) == 0 {
			return
		}
		if  err1 != nil {
			fmt.Println("@ Scan failed")
			fmt.Println(err1.Error())
		}
		ProcessDataSnapshot(schemaTarget, dbTarget, dbRoot, values)
		timekey := fmt.Sprintf("%s", values)
		sqlDelete := fmt.Sprintf("delete from snapshot where time = '%s';",timekey[:27])
		//fmt.Println(sqlDelete)
		dbRoot.Exec(sqlDelete)
	}

}

func ProcessDataSnapshot(schemaTarget *SchemaDb, dbTarget *sqlx.DB,dbRoot *sqlx.DB, values []interface{}) {
	nameTableInsert := fmt.Sprintf("%v",values[2])
	caseStatus := fmt.Sprintf("%s", values[3])
	if caseStatus == "delete" {
		fmt.Println("case delete not availability")
		return
	}
	allVal := convertInterfaceToString(fmt.Sprintf("%s", values[1]))
	allValues := []string{}
	allNameKeys := schemaTarget.ListKeys[nameTableInsert]
	allValuesNotKey := []string{}
	for i, v := range allVal{
		if i > len(allNameKeys) {
			allValuesNotKey = append(allValuesNotKey, schemaTarget.ListColumns[nameTableInsert][i] + " = " + fmt.Sprintf("'%s'", v))
		}
		allValues = append(allValues, fmt.Sprintf("'%s'", v))
	}

	sqlx := tpl.InsertTpl(nameTableInsert,allValues, allNameKeys, allValuesNotKey)
	_, err := dbTarget.Exec(sqlx)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(sqlx)
		return
	}
}

func convertInterfaceToString(value string) []string {
	value = value[1:len(value) -1 ]
	res := strings.Split(value,",")
	return res
}


func ProcessDataInTable( schemaRoot *SchemaDb, schemaTarget *SchemaDb, dbRoot *sqlx.DB, dbTarget *sqlx.DB, tableName string) {
	sqlX := tpl.SelectTpl(tableName,schemaRoot.ListColumns[tableName])
	rows, err := dbRoot.Queryx(sqlX)
	if err != nil {
		fmt.Println("@ query failed")
		fmt.Println(sqlX)
		fmt.Println(err.Error())
		return
	}

	for rows.Next() {
		values,err1 := rows.SliceScan()
		if  err1 != nil {
			fmt.Println("@ Scan failed")
			fmt.Println(err1.Error())

		}
		//fmt.Printf("%scheck interface\n", values)
		ProcessDataRow(schemaTarget, dbTarget,values,tableName,schemaRoot.ListNameTable)
		//break
	}

}

func ProcessDataRow(schemaTarget *SchemaDb, dbTarget *sqlx.DB, values []interface{}, tableName string, listTable []string) {
	index := -1
	for i, v := range listTable{
		if  v == tableName{
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("@ ProcessDataRow get index columns failed")
		return
	}
	nameTableInsert := schemaTarget.ListNameTable[index]
	allValues := []string{}
	allNameKeys := schemaTarget.ListKeys[nameTableInsert]
	allValuesNotKey := []string{}
	for i, v := range values{
		if i > len(allNameKeys) {
			allValuesNotKey = append(allValuesNotKey, schemaTarget.ListColumns[nameTableInsert][i] + " = " + fmt.Sprintf("'%s'", v))
		}
		allValues = append(allValues, fmt.Sprintf("'%s'", v))
	}
	sqlx := tpl.InsertTpl(nameTableInsert,allValues, allNameKeys, allValuesNotKey)
	_, err := dbTarget.Exec(sqlx)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(sqlx)
		return
	}
	//fmt.Println(sqlx)

}

func GetSchemaDatabase() (SchemaDb, error) {
	err := sqlDbInstance.GetListTable()
	if err != nil {
		return schemaDbInstance, err
	}
	schemaD := schemaDbInstance.ListNameTable
	for _, v := range schemaD {
		err1 := sqlDbInstance.GetListColumns(v)
		if err1 != nil {
			return schemaDbInstance, err1
		}
		err2 := sqlDbInstance.GetListKeys(v)
		if err2 != nil {
			return schemaDbInstance, err2
		}
	}
	return schemaDbInstance, nil
}