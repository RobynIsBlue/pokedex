package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func() error
}

type config struct {
	next     int
	previous int
}

var mappy = map[string]cliCommand{}

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
			if _, ok := mappy[cleanedInput[0]]; ok {
				mappy[cleanedInput[0]].callback()
			} else {
				fmt.Print("Unknown command\n")
			}
			fmt.Print("Pokedex > ")
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func byeBye() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func hewp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for key, val := range mappy {
		fmt.Printf("%s: %s\n", key, val.description)
	}
	return nil
}

type locationData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func mapFunc() error {
	for i := 0; i < 20; i++ {
		locDat, err := callApiWithID(mappy["map"].config.next)
		if err != nil {
			fmt.Print("getting location data\n")
			return err
		}
		fmt.Printf("%s\n", locDat.Name)
		mappy["map"].config.next++
	}
	return nil
}


func mapbFunc() error {
	if mappy["map"].config.next <= 20 {
		fmt.Print("you're on the first page!\n")
		return nil
	}
	mappy["map"].config.next -= 20
	mappy["map"].config.previous -= 20
	for i := 0; i < 20; i++ {
		locDat, err := callApiWithID(mappy["map"].config.next)
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


func callApiWithID(id int) (locationData, error) {
	link := "https://pokeapi.co/api/v2/location-area/" + fmt.Sprintf("%d", id) + "/"
	res, err := http.Get(link)
	if err != nil {
		fmt.Print("get link\n")
		return locationData{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print("error after reading\n")
		return locationData{}, err
	}
	res.Body.Close()
	var place locationData
	err = json.Unmarshal(body, &place)
	if res.StatusCode > 299 {
		fmt.Print("status code too high\n")
		return locationData{}, fmt.Errorf("%d", res.StatusCode)
	}
	if err != nil {
		fmt.Print("error after unmarshalling\n")
		return locationData{}, err
	}
	return place, nil
}
