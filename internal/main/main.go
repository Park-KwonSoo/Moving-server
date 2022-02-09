package main

import (
	"sync"

	Router "github.com/Park-Kwonsoo/moving-server/internal/router"
	CacheServer "github.com/Park-Kwonsoo/moving-server/pkg/cache-server"
	nosqlDB "github.com/Park-Kwonsoo/moving-server/pkg/database/nosql"
	_ "github.com/Park-Kwonsoo/moving-server/pkg/database/sql"
)

func main() {

	var wait sync.WaitGroup
	wait.Add(4)

	go func() {
		defer wait.Done()
		Router.SetupGRPCRouter()
	}()

	go func() {
		defer wait.Done()
		Router.SetupRESTRouter()
	}()

	go func() {
		defer wait.Done()
		CacheServer.Connect()
	}()

	go func() {
		defer wait.Done()
		nosqlDB.Connect()
	}()

	wait.Wait()
}
