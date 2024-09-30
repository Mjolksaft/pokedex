package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokedex/internal/cache"
	"pokedex/structs"
)

func GetLocations(currentLocation int) (structs.Locations, error) {

	var items structs.Locations
	data, exists := cache.GetChache(string(currentLocation))

	if exists { // if the data already exists we unmarshal it and send it by
		err := json.Unmarshal(data, &items)
		if err != nil {
			return structs.Locations{}, err
		}

	} else { // if it doesnt exist in chache we query for it

		fullUrl := "https://pokeapi.co/api/v2/location"

		//create a req
		res, err := http.Get(fullUrl)
		if err != nil {
			return structs.Locations{}, fmt.Errorf("error on req: %w", err)
		}
		defer res.Body.Close()

		// Decode the res
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&items); err != nil {
			return structs.Locations{}, fmt.Errorf("error on decoding: %w", err)
		}

		//save to chache or update it
		cache.Add(items, string(currentLocation))
	}
	return items, nil
}

func GetLocation(id string, currentLocation string) (structs.Location, error) {
	var item structs.Location
	fullKey := currentLocation + id

	data, exists := cache.GetChache(fullKey)

	if exists {
		err := json.Unmarshal(data, &item)
		if err != nil {
			return structs.Location{}, err
		}

	} else {
		fullUrl := "https://pokeapi.co/api/v2/location-area/" + id

		//create a req
		res, err := http.Get(fullUrl)
		if err != nil {
			return structs.Location{}, fmt.Errorf("error on req: %w", err)
		}
		defer res.Body.Close()

		// Decode the res
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&item); err != nil {
			return structs.Location{}, fmt.Errorf("error on decoding: %w", err)
		}
		cache.Add(item, string(fullKey))
	}

	return item, nil
}

func GetPokemonFromArea(id string, currentLocation string) (structs.PokemonEncounters, error) {
	var item structs.PokemonEncounters
	fullKey := currentLocation + id

	data, exists := cache.GetChache(fullKey)

	if exists {
		err := json.Unmarshal(data, &item)
		if err != nil {
			return structs.PokemonEncounters{}, err
		}

	} else {
		fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", id)

		//create a req
		res, err := http.Get(fullUrl)
		if err != nil {
			return structs.PokemonEncounters{}, fmt.Errorf("error on req: %w", err)
		}
		defer res.Body.Close()

		// Decode the res
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&item); err != nil {
			return structs.PokemonEncounters{}, fmt.Errorf("error on decoding: %w", err)
		}
		cache.Add(item, fullKey)
	}

	return item, nil
}

func GetPokemon(id string, currentLocation string) (structs.Pokemon, error) {
	var pokemon structs.Pokemon
	fullkey := currentLocation + id
	encoded, exists := cache.GetChache(fullkey)
	if exists {
		if err := json.Unmarshal(encoded, &pokemon); err != nil {
			return structs.Pokemon{}, err
		}

	} else {
		// make req
		fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
		res, err := http.Get(fullUrl)
		if err != nil {
			return structs.Pokemon{}, fmt.Errorf("no such pokemon %w", err)
		}

		decoder := json.NewDecoder(res.Body)
		// decode res
		if err := decoder.Decode(&pokemon); err != nil {
			return structs.Pokemon{}, fmt.Errorf("no such pokemon %w", err)
		}

		cache.Add(pokemon, fullkey)
	}

	return pokemon, nil
}
