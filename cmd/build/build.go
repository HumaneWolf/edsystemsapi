package main

import (
	"log"

	"humanewolf.com/ed/systemapi/systems"
)

func main() {
	log.Println("Welcome to the ED Systems API!")
	systems.BuildNameSearchTree()
}
