package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			GetDatabaseInstance()
		}()
	}

	wg.Wait()
}

type Database struct{}

var db *Database
var lock sync.Mutex

func (Database) SingletonInstance() *Database {
	fmt.Println("Creating Singleton Connection...\n")
	time.Sleep(2 * time.Second)
	fmt.Println("Singleton Connection Successfull...\n")

	return &Database{}
}

func GetDatabaseInstance() *Database {
	lock.Lock()
	defer lock.Unlock()

	if db == nil {
		println("Connecting database...\n")
		db = &Database{}

		db.SingletonInstance()

	} else {
		println("Database already Create...\n")
	}

	return db
}
