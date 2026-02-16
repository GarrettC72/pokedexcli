package pokeapi

type (
	PokeAPIResponse struct {
		Count    int      `json:"count"`
		Next     *string  `json:"next"`
		Previous *string  `json:"previous"`
		Results  []Result `json:"results"`
	}
	Result struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)
