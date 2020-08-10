package main

import (
	"demo/inventoryservice/product"
	"log"
	"net/http"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
