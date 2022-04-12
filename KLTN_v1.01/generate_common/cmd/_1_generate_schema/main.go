package main

import (
	"fmt"
	"generate_common/Utilities/common"
	"generate_common/Utilities/dbHandler"
	"generate_common/Utilities/fileHandler"
)

var (
	outputPath = "./configs/schema_output/"
	outputSchema = "schema.xlsx"
)

const (
	HostMaster       = "localhost"
	PortMaster       = 5432
	UserMaster       = "admin1"
	PasswordMaster   = "123456123"
	DbNameMaster     = "myapp"
	DriverNameMaster = "postgres"
	DriverNameTarget = "postgres"
)

func main() {
	fmt.Println("------------------------> start main ")
	clearFolder := common.ClearFolder(outputPath)
	if clearFolder != nil {
		return
	}
	psqlInfoMaster := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable ",
		HostMaster, UserMaster, PasswordMaster, PortMaster, DbNameMaster)

	err := dbHandler.InitDB(DriverNameMaster, psqlInfoMaster)
	if err != nil{
		fmt.Println("connect database failed")
		return
	}
	//dbHandler.Test()
	_ , o2 := dbHandler.GetSchemaDatabase()
	if o2 != nil {
		fmt.Println("get schema rdbRepository error")
		return
	}
	fileHandler.ProcessSetSchema(outputPath, outputSchema)

	//fileHandler.GetUpdate()

	fmt.Println("------------------------> success,  end main")
}


