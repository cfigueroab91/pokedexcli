package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cfigueroab91/pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

// getJSON revisa el caché y, si no hay datos, hace la petición.
// Solo guarda en el caché cuando el status fue correcto y el JSON es válido.
func (c Client) getJSON(url string, dest any) error {
	if data, ok := c.cache.Get(url); ok {
		if err := json.Unmarshal(data, dest); err == nil {
			return nil
		}
		// Bytes cacheados inválidos: lo tratamos como un miss.
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}

	c.cache.Add(url, data)
	return nil
}
