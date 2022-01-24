package main

import (
	"sync"

	Router "github.com/Park-Kwonsoo/moving-server/internal/router"
	nosqlDB "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
	_ "github.com/Park-Kwonsoo/moving-server/pkg/database/sql"
)

func main() {

	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()
		Router.SetupRouter()
	}()

	go func() {
		defer wait.Done()
		nosqlDB.Connect()
	}()

	wait.Wait()
}
