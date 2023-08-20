package main

import (
	"encoding/json"
	"fmt"
)

type Cat struct {
	// lowercase field cannot be exported
	// `json:"name"` will make Field Name become name in json string after apply json.Marshal()
	Name string `json:"name"`
	Age int
	IsAdult bool
}

func main() {
	data, _ := json.Marshal(Cat{
		Name: "Kitten",
		Age: 2,
		IsAdult: true,
	})

	cat := Cat{}
	 _ = json.Unmarshal(data, &cat)
	fmt.Println(cat)
}