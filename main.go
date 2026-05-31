package main

import (
	"fmt"
	"holiday/api"
)

func main() {

	holidays, err := api.GetHolidays()
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		return
	}

	for _, v := range holidays {
		fmt.Println(v)
	}
}
