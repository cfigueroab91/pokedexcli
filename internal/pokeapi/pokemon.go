package pokeapi

const PokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (c Client) GetPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon
	if err := c.getJSON(PokemonURL+name, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}
