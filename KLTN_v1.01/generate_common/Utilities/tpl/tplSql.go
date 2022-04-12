package tpl

import (
	"fmt"
	"strings"
)

func CreateTpl(nameTable string,columns []string, types[]string, keys []string) string {
	sqlA := fmt.Sprintf("create table if not exists %v (",nameTable)
	arrColumns := []string{}

	sqlZ :=	"\n);"
	for i, v := range columns{
		n := ""
		n += "\n\t" + v + " " + types[i]
		arrColumns = append(arrColumns, n)
	}
	sqlB := strings.Join(arrColumns, ",")
	pk := strings.Join(keys, ", ")
	if len(keys) != 0 {
		sqlB += ","
		pk = "\n\tprimary key (" + pk + ")"
	}
	return sqlA + sqlB + pk + sqlZ
}
//func DeleteTpl(nameTable string, key []string)  {
//	sqlx := fmt.Sprintf("delete from %v where ")
//
//}
func InsertTpl(nameTable string,arrValues []string,arrkeys []string,arrvalues []string) string {
	allValues := strings.Join(arrValues, ", ")
	keys := strings.Join(arrkeys, ", \n")
	values := strings.Join(arrvalues, ", \n")
	return "INSERT INTO " + nameTable +
		"\nvalues (" +
		"\n"+ allValues +
		"\n)" +
		"\non conflict (" + keys + ")" +
		"\ndo update set" +
		"\n" + values + " ;"
}
func SelectTpl(nameTable string,columns []string) string  {
	return "select * from "+ nameTable + ";"
	//return "SELECT ARRAY[ " +strings.Join(columns, "::text,")+ "]from "+nameTable +" ;"
}
