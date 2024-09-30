package menu

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokedex/internal"
	"strings"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func([]string) error
}

const (
	Mainloop = iota
	Maploop
)

var currentState int // keep track of the current loop we are in

var Menus []map[string]CliCommand

func init() {
	currentState = Mainloop
	Menus = []map[string]CliCommand{

		// for the main
		{
			"help": {
				Name:        "help",
				Description: "Displays a help message",
				Callback:    commandHelp,
			},
			"exit": {
				Name:        "help",
				Description: "exit the pokedex",
				Callback:    commandExit,
			},
			"map": {
				Name:        "help",
				Description: "Displays all locations",
				Callback:    commandMap,
			},
			"catch": {
				Name:        "catch",
				Description: "Catch a pokemon",
				Callback:    commandCatch,
			},
			"pokedex": {
				Name:        "pokedex",
				Description: "View pokedex",
				Callback:    commandPokedex,
			},
			"inspect": {
				Name:        "inspect",
				Description: "view stats of pokemon",
				Callback:    commandInspect,
			},
		},

		//for the map
		{
			"help": {
				Name:        "help",
				Description: "Displays a help message",
				Callback:    commandHelp,
			},
			"explore": {
				Name:        "explore",
				Description: "explore for pokemon 'explore id'",
				Callback:    commandExplore,
			},
			"back": {
				Name:        "back",
				Description: "back current menu",
				Callback:    commandBack,
			},
		},
	}
}

func commandHelp(commands []string) error {
	fmt.Println("Available commands")
	cnt := 1
	for key := range Menus[currentState] {
		fmt.Printf("%d. %s\n", cnt, key)
		cnt++
	}
	fmt.Println("")
	return nil
}

func commandExit(commands []string) error {
	fmt.Println("Exiting program")
	os.Exit(0)
	return nil
}

func commandMap(commands []string) error {
	items, err := internal.GetLocations(currentState)
	if err != nil {
		return fmt.Errorf("wrog with request: %w", err)
	}

	cnt := 1
	for _, items := range items.Results {
		fmt.Printf("%d. %s\n", cnt, items.Name)
		cnt++
	}

	ClILoop("map > ", Maploop)

	return nil
}

func commandExplore(commands []string) error {
	if len(commands) < 2 {
		return fmt.Errorf("not enough arguments")
	}
	location, err := internal.GetPokemonFromArea(commands[1], string(currentState))

	if err != nil {
		return fmt.Errorf("error %w", err)
	}

	for i, encounter := range location.PokemonEncounters {
		fmt.Printf("%d. %s\n", i+1, encounter.Pokemon.Name)
	}
	return nil
}

func commandBack(commands []string) error {
	currentState = Mainloop
	return nil
}

func commandCatch(commands []string) error {
	if len(commands) < 2 {
		return fmt.Errorf("wrong format must be 'catch name'")
	}

	// check if valid pokemon
	pokemon, err := internal.GetPokemon(commands[1], string(currentState))
	if err != nil {
		return err
	}
	number := rand.Intn(100)
	// get random number
	if number >= 50 {
		fmt.Printf("You caught: %s\n", pokemon.Name)
		err := internal.AddPokemon(commands[1]+string(currentState), pokemon)
		if err != nil {
			return fmt.Errorf("error adding pokemon to pokedex")
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandPokedex(commands []string) error {
	pokemons, err := internal.GetPokemonsFromStorage()
	if err != nil {
		return err
	}

	for i, pokemon := range pokemons {
		fmt.Printf("%d. %s\n", i+1, pokemon.Name)
	}

	return nil
}

func commandInspect(commands []string) error {
	if len(commands) < 2 {
		return fmt.Errorf("wrong format must be 'catch name'")
	}
	fullKey := commands[1] + string(currentState)

	data, exists := internal.GetPokemonFromPokedex(fullKey)

	if !exists {
		return fmt.Errorf("does not know pokemon")

	}

	fmt.Println(data)

	return nil
}

func GetCurrentState() int {
	return currentState
}

func SetCurrentState(state int) {
	currentState = state
}

func ClILoop(prefix string, menu int) {
	currentState = menu
	reader := bufio.NewReader(os.Stdin)
	for currentState == menu {
		fmt.Print(prefix)

		data, err := reader.ReadString('\n')
		cleanData := strings.TrimSpace(data)
		if err != nil {
			fmt.Println("error: %w", err)
		}

		splitData := strings.Split(cleanData, " ")
		command, exists := Menus[currentState][splitData[0]]

		if !exists {
			fmt.Println("No such command")
			continue
		}

		fmt.Println(command.Description)
		err = command.Callback(splitData)
		if err != nil {
			fmt.Println(err)
		}
	}
}
