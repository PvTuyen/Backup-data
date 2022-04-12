package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)


func ClearFolder(path string)  error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		e := os.Remove(path + f.Name())
		if e != nil {
			log.Fatal(e)
			return e
		}
	}
	return nil
}
func CheckKey(v1 string, keys []string) bool {
	for _, v:= range keys{
		if v1 == v {
			return true
		}
	}
	return false
}
func SaveFile(text string, path string, nameFile string)  {
	f ,_ := os.Create(path + nameFile)
	defer f.Close()
	_,  err := f.WriteString(text)
	if err != nil {
		fmt.Println("@write output file error")
		return
	}
}