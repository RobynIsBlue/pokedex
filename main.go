package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RobynIsBlue/pokedex/internal/pokecache"
)

const (
	intervalForCache = time.Second * 5
)

var cachePoke = pokecache.NewCache(intervalForCache)

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(i string) error
}

type config struct {
	next     int
	previous int
}

var mappy = map[string]cliCommand{}


func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}



func byeBye(i string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}



func hewp(i string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for key, val := range mappy {
		fmt.Printf("%s: %s\n", key, val.description)
	}
	// fmt.Println("\n")
	return nil
}


type pokemonEncounters struct {
	pokemonname string
}


type locationData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	PokemonEncounters []pokemonEncounters `json:"pokemon_encounters"`
}



func mapFunc(i string) error {
	for range 20 {
		locDat, err := callLocationApiWithID(mappy["map"].config.next)
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
		locDat, err := callLocationApiWithID(mappy["map"].config.next)
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


func explore(i string) error {
	if i == "" {
		fmt.Println("No location specified")
		return nil
	}
	// locDat, err := callLocationApiWithID(i)

	// if err != nil {
	// 	fmt.Println("Location not found")
	// 	return nil
	// }

	// pokemans := locDat.PokemonEncounters.Pokemon
	return nil
}


func callLocationApiWithID(id any) (locationData, error) {
	link := "https://pokeapi.co/api/v2/location-area/" + fmt.Sprintf("%v", id) + "/"

	body, ok := cachePoke.Get(link);
	if !ok {
		res, err := http.Get(link)
		if err != nil {
			fmt.Print("get link\n")
			return locationData{}, err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Print("error after reading\n")
			return locationData{}, err
		}
		res.Body.Close()

		if res.StatusCode > 299 {
			fmt.Print("status code too high\n")
			return locationData{}, fmt.Errorf("%d", res.StatusCode)
		}

		cachePoke.Add(link, body)
	}

	var place locationData
	err := json.Unmarshal(body, &place)
	
	if err != nil {
		fmt.Print("error after unmarshalling\n")
		return locationData{}, err
	}
	return place, nil
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
			name: "explore",
			description: "Explores pokemon in current area",
			config: &conf,
			callback: explore,
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
