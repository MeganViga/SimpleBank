package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	db "github.com/meganviga/simple_bank/db/sqlc"
	"github.com/meganviga/simple_bank/util"
	"github.com/meganviga/simple_bank/api"
)
/*const (
	dbDriver = "postgres"
	dbSource ="postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	sourceAddress ="0.0.0.0:8082"
)*/
func main(){
	config, err := util.LoadConfig(".")
	fmt.Println("configs",config)
	if err != nil{
		log.Fatal("Cannot load config file", err)
	}
	conn, err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil{
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(config.ServerAddress)
}