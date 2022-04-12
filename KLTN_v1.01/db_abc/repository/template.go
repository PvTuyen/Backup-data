package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)
var temp *template.Template

func main() {
	temp = template.Must(template.ParseFiles("template-06.txt"))
	fmt.Println("main")
	cuteAnimalsSpecies := map[string]string{
		"Dogs": "German Shepherd",
		"Cats": "Ragdoll",
		"Mice": "Deer Mouse",
		"Fish": "Goldfish",
	}
	err := temp.Execute(os.Stdout, cuteAnimalsSpecies)
	if err != nil {
		fmt.Println("error")
		log.Fatalln(err)
	}
}