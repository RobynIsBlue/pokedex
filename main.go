package main

import (
	"bufio"
	"fmt"
	"os"
)

var mappy = map[string]cliCommand{}

func hewp(_ string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for key, val := range mappy {
		fmt.Printf("%s: %s\n", key, val.description)
	}
	// fmt.Println("\n")
	return nil
}

func mapFunc(i string) error {
	for range 20 {
		locDat, err := callLocationApi(mappy["map"].config.next)
		if err != nil {
			fmt.Print("getting location data\n")
			return err
		}
		fmt.Printf("%s\n", locDat.Name)
		mappy["map"].config.next++
	}
	return nil
}

func mapbFunc(i string) error {
	if mappy["map"].config.next <= 20 {
		fmt.Print("you're on the first page!\n")
		return nil
	}

	mappy["map"].config.next -= 20
	mappy["map"].config.previous -= 20

	for range 20 {
		locDat, err := callLocationApi(mappy["map"].config.next)
		if err != nil {
			fmt.Print("getting location data\n")
			return err
		}
		fmt.Printf("%s\n", locDat.Name)
		mappy["map"].config.next++
	}
	mappy["map"].config.next -= 20
	mappy["map"].config.previous -= 20
	return nil
}

func main() {

	conf := config{
		next:     1,
		previous: -1,
	}

	mappy = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			config:      &conf,
			callback:    byeBye,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			config:      &conf,
			callback:    hewp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			config:      &conf,
			callback:    mapFunc,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			config:      &conf,
			callback:    mapbFunc,
		},
		"explore": {
			name:        "explore",
			description: "Explores pokemon in current area",
			config:      &conf,
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			config:      &conf,
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a captured pokemon",
			config:      &conf,
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Check caught pokemon in your pokedex",
			config:      &conf,
			callback:    pokerdex,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			text := scanner.Text()
			cleanedInput := cleanInput(text)
			if len(cleanedInput) == 0 {
				fmt.Print("Pokedex > ")
				continue
			}

			if k, ok := mappy[cleanedInput[0]]; ok {
				parameter := ""
				if len(cleanedInput) > 1 {
					parameter = cleanedInput[1]
				}
				k.callback(parameter)
			} else {
				fmt.Print("Unknown command\n")
			}
			fmt.Print("Pokedex > ")
		}
	}
}
