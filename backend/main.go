package main

import (
	"github.com/git-amw/backend/routers"
)

func main() {
	route := routers.SetupRouter()
	route.Run("localhost:8000")
}
