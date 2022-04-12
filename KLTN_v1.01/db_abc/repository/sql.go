package main

import (
	"fmt"
	"github.com/huandu/go-sqlbuilder"
)

func main()  {
	sql := sqlbuilder.Select("id", "name").From("demo.user").
		Where("status = 1").Limit(10).
		String()

	fmt.Println(sql)
}

