package main

import (
	"fmt"
	"generate_common/Utilities/dbHandler"
	"generate_common/Utilities/fileHandler"
)

const (
	HostRoot       = "172.16.210.31"
	PortRoot       = 5432
	UserRoot       = "admin"
	PasswordRoot   = "123456123"
	DbNameRoot     = "kltn"
	DriverNameRoot = "postgres"
)
const (
	HostTarget       = "172.16.210.31"
	PortTarget       = 5432
	UserTarget       = "admin"
	PasswordTarget   = "123456123"
	DbNameTarget     = "db_target"
	DriverNameTarget = "postgres"
)
var path = "./configs/schema_output/"
var schemaInput = "schema.xlsx"
func main() {
	fmt.Println("------------------------> start main ")
	psqlInfoRoot := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable ",
		HostRoot, UserRoot, PasswordRoot, PortRoot, DbNameRoot)
	dbRoot, err0 := dbHandler.ConnectDb(DriverNameRoot, psqlInfoRoot)
	if err0 != nil {
		fmt.Println(err0.Error())
		return
	}
	psqlInfoTarget := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable ",
		HostTarget, UserTarget, PasswordTarget, PortTarget, DbNameTarget)
	dbTarget, err1 := dbHandler.ConnectDb(DriverNameTarget, psqlInfoTarget)
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	schemaRoot := fileHandler.GetSchemaDesign(path,schemaInput,"root")
	//fmt.Println(schemaRoot.ListNameTable)
	schemaTarget := fileHandler.GetSchemaDesign(path,schemaInput,"target")
	//fmt.Println(*schemaTarget)
	dbHandler.ProcessSchema(schemaRoot,schemaTarget,dbRoot,dbTarget)
	for  {
		dbHandler.UpdaterSql(schemaRoot,schemaTarget,dbRoot,dbTarget)
	}
	fmt.Println("------------------------> success,  end main")

}