package steamid

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type response struct {
	Data data `json:"response"`
}

type data struct {
	SteamID64 string `json:"steamid"`
	Success   int    `json:"success"`
}

func ConvertVanityURL(vkey, vanityurl string) (string, error) {
	u, err := url.Parse("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/")
	if err != nil {
		return "", err
	}

	// build our request
	q := u.Query()
	q.Set("key", vkey)
	q.Set("vanityurl", vanityurl)
	u.RawQuery = q.Encode()

	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	response := response{}

	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &response)

	if response.Data.Success != 1 {
		return "", errors.New("could not convert URL")
	}

	return response.Data.SteamID64, nil
}
