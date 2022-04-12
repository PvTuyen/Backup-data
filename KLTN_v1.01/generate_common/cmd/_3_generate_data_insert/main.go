package main

import (
	"fmt"
	"generate_common/Utilities/dbHandler"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	dbTarget *sqlx.DB
	inputPath = "./configs/sql_table_output/"
	outputPath = "./configs/data_insert_output/"
	sizeInsert = 100
	lengthText = 20
	insertSql = ""
	dataText = "qwercvbnm123asdfgbnm123asdfgh1jtcvbnh1jtcbnm123asdfgh1jtcvbnvbbnm123asdfgh1jtcvbnnm123bnm123asdfgh1jtcvbnasdfghj1zxcvbxcvb1nm1xcvbnm1xcvbnm12cvbnm123asdfghj22nm123yuiopzx90QWERTcvbnm123asdfghjYUbnm123asdfgh1jtcvbnzxcvbnm123IOPASDFGHJzxcvbnm123Kklzxccvbnm123asdfghjvbnm1234567890QWERTYUzxcvbnm123IOcvbnm123asdfghjPASDFGHJzxcvbnm123cvbnm123asdfghjKLZXCVBzxcvbnm123NM"
)

const (
	Host       = "172.16.210.31"
	Port       = 5432
	User       = "admin"
	Password   = "123456123"
	DbName     = "kltn"
	DriverName = "postgres"
)

func main() {
	psqlInfoTarget := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable ",
		Host, User, Password, Port, DbName)
	db, err1 := dbHandler.ConnectDb(DriverName, psqlInfoTarget)
	dbTarget = db
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	rand.NewSource(time.Now().UnixNano())
	listCreateTable := getAllTable(inputPath)
	for _, v := range listCreateTable{
		_, err := dbTarget.Exec(v)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(v)
			continue
		}
		generateDataInsert(v)
	}
}

func generateDataInsert(createSql string) {
	//createSql = strings.ReplaceAll(createSql, "  ", " ")
	allCreateTable := strings.Split(createSql, ");")
	//fmt.Println(createSql)
	for _, one := range allCreateTable{
		name, insert := genInsert(one,sizeInsert)
		if name == "" {
			continue
		}
		_, err := dbTarget.Exec(insert)
		if err != nil {
			fmt.Println(err.Error())
		}
		//common.SaveFile(insert,outputPath, name+".sql")
	}

}

func genInsert(one string, size int) (string,string) {
	rand.NewSource(time.Now().UnixNano())
	if !strings.Contains(one,"create table if not exists") {
		return "",""
	}
	//fmt.Println("\n check",one)
	index := 12
	nameTblBf := strings.Index(one, "create table if not exists")
	nameTblAf := strings.Index(one, "(")

	nameTable :=  one[nameTblBf + len("create table if not exists"):nameTblAf]

	values := strings.Split(one[index:], ",")
	arrayRes := []string{}

	for i := 0; i < size; i++ {
		arrValues := []string{}
		for _, v := range values{
			if strings.Contains(v, "key") {
				break
			}
			var dataS = ""
			var dataN float32
			if strings.Contains(v,"text") {
				//dataS = "A"
				dataS = "'" +getText() +"'"
			} else if strings.Contains(v,"numeric") {
				dataN = rand.Float32() / rand.Float32()
			} else {
				dataS = "null"
			}

			if dataN  !=  0{
				arrValues = append(arrValues, fmt.Sprintf("%v",dataN))
			} else {
				//fmt.Println(dataN)
				arrValues = append(arrValues, fmt.Sprintf("%s",dataS))
			}
		}
		insert := "insert into " + nameTable + " values (" + strings.Join(arrValues, ", ") + ");"
		arrayRes = append(arrayRes, insert)
	}
	res := strings.Join(arrayRes, "\n")
	return nameTable,res
}

func getText() string {
	res := ""
	for i := 0; i < lengthText; i++ {
		index := rand.Intn(len(dataText) - 1)
		res = fmt.Sprintf("%v%v", res, string(dataText[index]))
	}
	return res
}

func getAllTable(path string) []string {
	arrSql := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	for _, file := range files {
		data , err1 := os.ReadFile(path+file.Name())
		if err1 != nil {
			return nil
		}
		sql := string(data)
		if sql == "" {
			fmt.Println("sql in file = null", path + file.Name())
			continue
		}
		arrSql = append(arrSql,sql)
	}
	return arrSql
}