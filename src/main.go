package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"demo/inventoryservice/database"
	"demo/inventoryservice/product"
	"demo/inventoryservice/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	receipt.SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
