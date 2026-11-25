package pokeapi

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

type LocationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationArea struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c Client) ListLocationAreas(url string) (LocationAreasResp, error) {
	var areas LocationAreasResp
	if err := c.getJSON(url, &areas); err != nil {
		return LocationAreasResp{}, err
	}
	return areas, nil
}

func (c Client) GetLocationArea(name string) (LocationArea, error) {
	var area LocationArea
	if err := c.getJSON(LocationAreaURL+name, &area); err != nil {
		return LocationArea{}, err
	}
	return area, nil
}
