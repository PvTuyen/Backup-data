package main

import (
	"fmt"
	"generate_common/Utilities/common"
	"generate_common/Utilities/fileHandler"
)

var (
	inputPath = "./configs/schema_output/"
	inputFile = "schema.xlsx"
	outputPath = "./configs/sql_table_output/"
	fileOut = "Create_table.sql"
)


func main() {

	fmt.Println("------------------------> start main ")
	clearFolder := common.ClearFolder(outputPath)
	if clearFolder != nil {
		return
	}
	ok := fileHandler.ProcessGetSchema(inputPath, inputFile, outputPath,fileOut)
	if ok != nil {
		fmt.Println("get schema in excel error")
		return
	}

	fmt.Println("------------------------> success,  end main")
}


