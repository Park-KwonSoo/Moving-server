package main

import (
	"sync"

	db "github.com/Park-Kwonsoo/moving-server/models"
	Router "github.com/Park-Kwonsoo/moving-server/router"
)

func main() {

	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer func() {
			db.Disconnect()
			wait.Done()
		}()
		db.Connect()
		Router.SetupRouter()
	}()

	wait.Wait()
}
