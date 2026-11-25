package main

import "fmt"

func commandMap(cfg *config, args []string) error {
	if cfg.nextLocationsURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return listLocations(cfg, *cfg.nextLocationsURL)
}

func commandMapb(cfg *config, args []string) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return listLocations(cfg, *cfg.prevLocationsURL)
}

func listLocations(cfg *config, url string) error {
	areas, err := cfg.pokeapiClient.ListLocationAreas(url)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = areas.Next
	cfg.prevLocationsURL = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}
