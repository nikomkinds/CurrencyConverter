package main

import (
	"github.com/nikomkinds/CurrencyConverter/internal/server"
	"log"
)

func main() {
	serv := server.NewServer()
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
	/*date := "01/01/2001"
	slice := date[6:]
	fmt.Println(slice)*/
}
