package fileHandler

import (
	"errors"
	"fmt"
	"generate_common/Utilities/common"
	"generate_common/Utilities/dbHandler"
	"generate_common/Utilities/tpl"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
	"strings"
	"time"
)

type excelMapSchema struct {
	name     string
	schemaDb *dbHandler.SchemaDb
}
var schema = excelMapSchema{
	schemaDb: dbHandler.GetSchemaInstance(),
}
var sqlCreate = []string{}

func ProcessGetSchema(path string, fileName string, outputPath string, fileOut string) error {
	err0 := errors.New("@ get schema failed  ")
	file , err := excelize.OpenFile(path+fileName)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("open file excel failed")
		return err0
	}
	sheetNames := file.GetSheetList()
	for _, v:= range sheetNames{
		err1 := getValueSheet(v,file)
		if err1 != nil {
			return err0
		}
	}
	for _, v := range schema.schemaDb.ListNameTable{
		sql := tpl.CreateTpl(v,schema.schemaDb.ListColumns[v],schema.schemaDb.ListType[v],schema.schemaDb.ListKeys[v])
		sqlCreate = append(sqlCreate,sql)
	}
	allSql := strings.Join(sqlCreate,"\n--------------------------------------------\n")

	common.SaveFile(allSql,outputPath,fileOut)
	return nil
}

func getValueSheet(nameSheet string, file *excelize.File) error {
	err := errors.New("@ error get value in sheet "+ nameSheet)
	nameTable, err0 := file.GetCellValue(nameSheet,"C4")
	nameTable = strings.TrimSpace(nameTable)
	if nameTable == "" {
		return nil
	}
	listNameCol := []string{}

	listKeys := []string{}
	listType := []string{}
	indexCol := 12
	for  {
		nameCol, err1 := file.GetCellValue(nameSheet,"D" + strconv.Itoa(indexCol))

		typeCol, err2 := file.GetCellValue(nameSheet,"H" + strconv.Itoa(indexCol))

		keys, err3 := file.GetCellValue(nameSheet,"G" + strconv.Itoa(indexCol))
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
			return err
		}

		nameCol = strings.TrimSpace(nameCol)
		typeCol = strings.TrimSpace(typeCol)
		keys = strings.TrimSpace(keys)
		if nameCol == "" || typeCol == "" {
			break
		}
		listNameCol = append(listNameCol, nameCol)
		listType = append(listType, typeCol)
		if keys != "" {
			listKeys = append(listKeys, nameCol)
		}
		indexCol++
	}
	schema.schemaDb.ListNameTable = append(schema.schemaDb.ListNameTable , nameTable)
	schema.schemaDb.ListColumns[nameTable] = listNameCol
	schema.schemaDb.ListType[nameTable] = listType
	schema.schemaDb.ListKeys[nameTable] = listKeys
	return nil
}
func GetUpdate()  {
	schema.schemaDb = dbHandler.GetSchemaInstance()
	listTable := schema.schemaDb.ListNameTable
	arrUp := []string{}
	arrDe := []string{}
	fmt.Println("size table= ", len(listTable))
	for _, v := range listTable{
		//fmt.Println(v)
		update := sqlUpdate(v,schema.schemaDb)
		delet := sqlDelete(v, schema.schemaDb)
		arrUp = append(arrUp, update)
		arrDe = append(arrDe, delet)
		//break
	}
	//fmt.Println(strings.Join(arrDe,"\n"))
	//fmt.Println(strings.Join(arrUp,"\n"))
	common.SaveFile(strings.Join(arrDe,"\n"),"./configs/data_insert_output/update_1k5.sql","")
	common.SaveFile(strings.Join(arrUp,"\n"),"./configs/data_insert_output/delete_1k5.sql","")

}

func sqlDelete(v string, db *dbHandler.SchemaDb) string {
	//sql := "set k0 db01='9999' where k001=(select tb.k001  from togo.unlisted tb order by tb.k001 limit 1);"
	columns := db.ListColumns[v]
	keys := db.ListKeys[v]
	strKey := strings.Join(keys,",")
	list := []string{}
	for i,_:= range  columns {
		if !strings.Contains(strKey, columns[len(columns) -1 - i]) {
			list = append(list, columns[len(columns) -1 - i])
			if len(list) == 2 {
				break
			}
		} else {
			break
		}
	}
	value := strings.Join(list,"= NULL, ")
	sqlx := fmt.Sprintf(" update togo.%v set %v= NULL where true ;",v,value)
	return sqlx
}

func sqlUpdate(v string, db *dbHandler.SchemaDb) string {

	//delete from public.nkstock02 r where r in (select t from public.nkstock02 t order by (k001,k002,k003) limit 1);"
	keys := db.ListKeys[v]
	strKey := strings.Join(keys,",")
	sql := fmt.Sprintf("delete from togo.%v r where r in (select t from togo.%v t order by (%v) limit 1);",v,v,strKey)
	return sql
}
func ProcessSetSchema(path string, nameFile string)  {
	schema.name = nameFile
	schema.schemaDb = dbHandler.GetSchemaInstance()
	f := excelize.NewFile()
	printfSheetInfo(f)
	for _, v := range schema.schemaDb.ListNameTable {
		if strings.Contains(v, "_snapshot") {
			continue
		}
		printfSheet(v , f)
	}
	f.DeleteSheet("Sheet1")
	if err := f.SaveAs(path +nameFile); err != nil {
		fmt.Println(err)
	}
}

func printfSheetInfo(f *excelize.File) {
	v := "info_database"
	index := f.NewSheet(v)
	printfSheetInfoBackground(v,f)
	f.SetActiveSheet(index)
}

func printfSheet( v string, f *excelize.File) {
	index := f.NewSheet(v)
	printfBackground(v,f)
	f.SetCellValue(v,"C3", v)
	f.SetCellValue(v,"C4", v)
	f.SetCellValue(v,"C5", fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05")))
	columns := schema.schemaDb.ListColumns[v]
	keys := schema.schemaDb.ListKeys[v]
	types := schema.schemaDb.ListType[v]
	countKey := 0
	for i, v1 := range columns{
		stt := i + 1
		f.SetCellValue(v,fmt.Sprintf("B%v", 12 + i), fmt.Sprintf("%v",stt))
		f.SetCellValue(v,fmt.Sprintf("C%v", 12 + i), v1)
		f.SetCellValue(v,fmt.Sprintf("D%v", 12 + i), v1)
		f.SetCellValue(v,fmt.Sprintf("F%v", 12 + i), types[i])
		f.SetCellValue(v,fmt.Sprintf("H%v", 12 + i), getNewType(types[i]))
		if common.CheckKey(v1,keys) {
			countKey++
			f.SetCellValue(v,fmt.Sprintf("E%v", 12 + i), countKey)
			f.SetCellValue(v,fmt.Sprintf("G%v", 12 + i), countKey)
		}
	}
	f.SetActiveSheet(index)
}

func printfSheetInfoBackground(v string, f *excelize.File) {
	style, _ := f.NewStyle(`{
	"font": {
		"bold": true,
		"italic": true,
		"family": "Times New Roman",
		"size": 36,
		"color": "#777777"
	},
	"fill": {
		"type": "pattern",
		
		"pattern": 1
	}
}`)
	f.SetColWidth(v, "B", "B", 30)
	f.SetColWidth(v, "A", "A", 50)
	f.SetCellValue(v, "A1", "Database root")
	f.SetCellValue(v, "B2", "Database type")
	f.SetCellValue(v, "B3", "IP address")
	f.SetCellValue(v, "B4", "Port")
	f.SetCellValue(v, "B5", "Database name")
	f.SetCellValue(v, "B6", "User")
	f.SetCellValue(v, "B7", "Password")

	f.SetCellValue(v, "A9", "Database target")
	f.SetCellValue(v, "B10", "Database type")
	f.SetCellValue(v, "B11", "IP address")
	f.SetCellValue(v, "B12", "Port")
	f.SetCellValue(v, "B13", "Database name")
	f.SetCellValue(v, "B14", "User")
	f.SetCellValue(v, "B15", "Password")
	style2, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#b1cefc"],"pattern":1},"font":{"bold":true,"size":13}}`)
	_ = f.SetCellStyle(v, "A1", "A1", style)
	_ = f.SetCellStyle(v, "A9", "A9", style)
	_ = f.SetCellStyle(v, "B2", "B7", style2)
	_ = f.SetCellStyle(v, "B10", "B15", style2)
}

func getNewType(typeColumn string) string {
	if dbHandler.DriverNameTarget == dbHandler.DriverNameTarget {
		return typeColumn
	}
	return "aaa"
}



func printfBackground(v string, f *excelize.File) {
	style, _ := f.NewStyle(`{
	"font": {
		"bold": true,
		"italic": true,
		"family": "Times New Roman",
		"size": 36,
		"color": "#777777"
	},
	"fill": {
		"type": "pattern",
		
		"pattern": 1
	}
}`)
	f.SetColWidth(v, "B", "B", 20)
	f.SetColWidth(v, "C", "D", 20)
	f.SetColWidth(v, "E", "F", 20)
	f.SetColWidth(v, "G", "H", 20)
	f.SetColWidth(v, "H", "K", 20)
	f.SetColWidth(v, "I", "I", 20)
	f.SetCellValue(v, "A1", "Overview")
	f.SetCellValue(v, "B3", "Root name table")
	f.SetCellValue(v, "B4", "New name table")
	f.SetCellValue(v, "B5", "Date")
	f.SetCellValue(v, "B6", "Comment")
	f.SetCellValue(v, "A9", "Columns")
	f.SetCellValue(v, "B11", "ON")
	f.SetCellValue(v, "C11", "Root name")
	f.SetCellValue(v, "D11", "New name")
	f.SetCellValue(v, "E11", "Root primary key ")
	f.SetCellValue(v, "F11", "Root Data type")
	f.SetCellValue(v, "G11", "New Primary key")
	f.SetCellValue(v, "H11", "New Data type")
	f.SetCellValue(v, "I11", "Comment")
	style2, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#b1cefc"],"pattern":1},"font":{"bold":true,"size":13}}`)
	style3, _ := f.NewStyle(`{"fill":{"type":"pattern","pattern":1},"font":{"bold":true,"size":13}}`)
	_ = f.SetCellStyle(v, "A1", "A1", style)
	_ = f.SetCellStyle(v, "A9", "A9", style)
	_ = f.SetCellStyle(v, "B3", "B7", style2)
	_ = f.SetCellStyle(v, "B11", "I11", style3)
}



func Exam(path string,  nameFile string)  {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	fmt.Println(index)
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(path +"Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func GetSchemaDesign(path string, nameFile string, dbName string) *dbHandler.SchemaDb {
	schemaDb := &dbHandler.SchemaDb{
		ListNameTable: []string{},
		ListColumns:   make(map[string][]string),
		ListKeys:      make(map[string][]string),
		ListType:      make(map[string][]string),
	}
	if dbName ==  "root" {
		return getSchema(schemaDb,"C3", "C", "E", "F", path + nameFile)
	}
	if dbName == "target" {
		return getSchema(schemaDb,"C4", "D", "E", "F", path + nameFile)
	}
	return schemaDb
}

func getSchema(db *dbHandler.SchemaDb,nameTable,  nameCol string, keyCol string, typeCol string, filePath string) *dbHandler.SchemaDb {
	file , err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("open file excel failed")
		return db
	}
	sheetNames := file.GetSheetList()
	for _, v:= range sheetNames{
		getSchemaTable(v, file,nameTable,db ,nameCol, keyCol, typeCol )

	}
	return db
}

func getSchemaTable(nameSheet string, file *excelize.File, tableName string, db *dbHandler.SchemaDb, ColIndex string, keyIndex string, typeIndex string) {
	nameTable, err0 := file.GetCellValue(nameSheet,tableName)
	nameTable = strings.TrimSpace(nameTable)
	if nameTable == "" {
		return
	}
	listNameCol := []string{}
	listKeys := []string{}
	listType := []string{}
	indexCol := 12
	for  {
		nameCol, err1 := file.GetCellValue(nameSheet,ColIndex + strconv.Itoa(indexCol))

		typeCol, err2 := file.GetCellValue(nameSheet,typeIndex + strconv.Itoa(indexCol))

		keys, err3 := file.GetCellValue(nameSheet,keyIndex + strconv.Itoa(indexCol))
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
			return
		}

		nameCol = strings.TrimSpace(nameCol)
		typeCol = strings.TrimSpace(typeCol)
		keys = strings.TrimSpace(keys)
		if nameCol == "" || typeCol == "" {
			break
		}
		listNameCol = append(listNameCol, nameCol)
		listType = append(listType, typeCol)
		if keys != "" {
			listKeys = append(listKeys, nameCol)
		}
		indexCol++
	}
	db.ListNameTable = append(db.ListNameTable , nameTable)
	db.ListColumns[nameTable] = listNameCol
	db.ListType[nameTable] = listType
	db.ListKeys[nameTable] = listKeys
}