package main

import (
	"fmt"
	"pokedex/internal/cache"
	"pokedex/internal/menu"
)

func main() {
	cache.StartInterval()

	fmt.Println("Welcome to the pokedex")
	fmt.Println("")

	menus := menu.Menus[0]
	menus["help"].Callback(nil)

	menu.ClILoop("pokedex > ", 0)
}
