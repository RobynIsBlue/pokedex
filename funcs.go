package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func byeBye(_ string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func explore(i string) error {
	if i == "" {
		fmt.Println("No location specified")
		return nil
	}
	locDat, err := callLocationApi(i)

	if err != nil {
		fmt.Println("Location not found")
		return nil
	}

	for _, v := range locDat.PokemonEncounters {
		fmt.Printf("- %s\n", v.Pokemon.PokemonName)
	}
	return nil
}

var pokedex = make(map[string]Pokemon)

func catch(i string) error {
	if i == "" {
		fmt.Println("No pokemon specified")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", i)
	pokeDat, err := callPokemonApi(i)
	if err != nil {
		return err
	}

	xp := pokeDat.BaseExperience
	switch {
	case xp <= rand.Intn(500):
		fmt.Printf("%s was caught!\n", i)
		pokedex[i] = pokeDat
	default:
		fmt.Printf("%s escaped!\n", i)
	}
	return nil
}

func inspect(i string) error {
	if i == "" {
		fmt.Println("No pokemon specified")
		return nil
	}
	if _, ok := pokedex[i]; !ok {
		fmt.Print("you have not caught that pokemon\n")
	}
	fmt.Printf("Height: %v\nWeight:%v\n", pokedex[i].Height, pokedex[i].Weight)
	pokedex[i].printStatsAndTypes()
	return nil
}

func pokerdex(_ string) error {
	if len(pokedex) == 0 {
		fmt.Println("No pokemon caught!")
	}
	for k := range pokedex {
		fmt.Printf("\t- %v\n", k)
	}
	return nil
}
