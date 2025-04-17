package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type Stat struct {
	Name string `json:"name"`
}

type Stats struct {
	Stat Stat `json:"stat"`
	BaseStat int  `json:"base_stat"`
}

// type Stats struct {
// 	Stat     struct {
// 		Name string `json:"name"`
// 	} `json:"stat"`
// }

type Type struct {
	Name string `json:"name"`
}

type Types struct {
	Type Type `json:"type"`
}

// do I need lists of the types here?
type Pokemon struct {
	PokemonName    string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type pokemonEncounters struct {
	Pokemon Pokemon
}

type locationData struct {
	ID                int                 `json:"id"`
	Name              string              `json:"name"`
	PokemonEncounters []pokemonEncounters `json:"pokemon_encounters"`
}

func (p Pokemon) printStatsAndTypes() {
	fmt.Println("Stats:")
	for _, v := range p.Stats {
		fmt.Printf("\t- %s: %v\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range p.Types {
		fmt.Printf("\t- %v\n", v.Type.Name)
	}
}

func callLocationApi(id any) (locationData, error) {

	link := "https://pokeapi.co/api/v2/location-area/" + fmt.Sprintf("%v", id) + "/"

	body, ok := cachePoke.Get(link)
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


func callPokemonApi(name string) (Pokemon, error) {
	link := "https://pokeapi.co/api/v2/pokemon/" + name + "/"
	
	res, err := http.Get(link)
	if err != nil {
		fmt.Println("getting pokemon info with link")
		return Pokemon{}, nil
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Print("status code too high\n")
		return Pokemon{}, fmt.Errorf("%d", res.StatusCode)
	}

	var poke Pokemon
	err = json.Unmarshal(body, &poke)
	if err != nil {
		fmt.Print("unmarshaling pokemon")
		return Pokemon{}, nil
	}
	return poke, nil
}

