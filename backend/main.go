package main

import (
	"github.com/git-amw/backend/db"
	"github.com/git-amw/backend/routers"
)

func main() {
	routers.DBInstance = db.Connect()
	routers.ESClient = db.ESClientConnection()
	// db.ESCheackIndexExists(routers.ESClient)
	route := routers.SetupRouter()
	route.Run("localhost:8000")
}
