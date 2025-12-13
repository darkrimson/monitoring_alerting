package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Worker started")

	for {
		// позже: выбор мониторов, выполнение проверок, логика инцидентов
		log.Println("Worker tick")
		time.Sleep(time.Second * 5)
	}
}
