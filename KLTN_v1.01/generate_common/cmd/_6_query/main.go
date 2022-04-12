package main

import (
	"fmt"
	"generate_common/Utilities/dbHandler"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	dbTarget *sqlx.DB
	inputPath = "./configs/sql_table_output/"
	)

const (
	HostTarget       = "172.16.210.31"
	PortTarget       = 5432
	UserTarget       = "admin"
	PasswordTarget   = "123456123"
	DbNameTarget     = "test_sv"
	DriverNameTarget = "postgres"
)

func main() {
	psqlInfoTarget := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable ",
		HostTarget, UserTarget, PasswordTarget, PortTarget, DbNameTarget)
	db, err1 := dbHandler.ConnectDb(DriverNameTarget, psqlInfoTarget)
	dbTarget = db
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	rand.NewSource(time.Now().UnixNano())
	proc(inputPath, db)
	fmt.Println("end main")
}

func queryTable(v string, db *sqlx.DB) {

}

func proc(path string, db *sqlx.DB) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)

	}

	for _, file := range files {
		data , err1 := os.ReadFile(path+file.Name())
		if err1 != nil {
			fmt.Println("@@@ read file failed ", path+file.Name())
			continue
		}
		sql := string(data)
		if sql == "" {
			fmt.Println("sql in file = null", path + file.Name())
			continue
		}
		//arrSql := strings.Split(sql,");")
		_, err2 := db.Exec(sql)
		if err2 != nil {
			fmt.Println("@@@ query failed ", path+file.Name())
			fmt.Println(err2)

		}

	}

}