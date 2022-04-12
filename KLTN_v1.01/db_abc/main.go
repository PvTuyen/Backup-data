package main

import (
	"fmt"
	"text/template"
)
var temp *template.Template

func main() {
	var abc int64 = 124
	var abc1 float64 = 1234.124
	if 1 < abc1 {
		fmt.Println("so sanh duoc")
	}
	fmt.Println("abc ", abc)
	fmt.Println("abc1", abc1)
}