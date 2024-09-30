package internal

import (
	"pokedex/structs"
	"sync"
)

type safePokedex struct {
	mu      *sync.RWMutex
	pokedex map[string]structs.Pokemon
}

var storage = safePokedex{mu: &sync.RWMutex{}, pokedex: make(map[string]structs.Pokemon)}

func AddPokemon(key string, data structs.Pokemon) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.pokedex[key] = data

	return nil
}

func GetPokemonsFromStorage() ([]structs.Pokemon, error) {
	var items []structs.Pokemon

	for _, pokemon := range storage.pokedex {
		items = append(items, pokemon)
	}

	return items, nil
}

func GetPokemonFromPokedex(key string) (structs.Pokemon, bool) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	value, exists := storage.pokedex[key]
	if !exists {
		return structs.Pokemon{}, exists
	}

	// fmt.Println("EXISTs IN CHACHE")
	return value, exists
}
