package main

import (
	"humanewolf.com/ed/systemapi/api"
	"humanewolf.com/ed/systemapi/systems"
)

func main() {
	systems.BuildNameSearchTree()
	api.RunAPI()
}
