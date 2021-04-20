package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	steamid "github.com/01dog/pcheck/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type response struct {
	Converted converted `json:"converted"`
}

type converted struct {
	SteamID64 string `json:"steamid64"`
	SteamID32 string `json:"steamid"`
	Steam3    string `json:"steam3"`
}

var (
	// used for flags
	skey     string
	vkey     string
	steamurl string
)

// idCmd represents the id command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "converts steam URLs into the 3 versions of steam IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		skey = viper.GetString("skey")
		vkey = viper.GetString("vkey")

		u, err := url.Parse(steamurl)
		if err != nil {
			fmt.Println("error parsing steam url: ", err)
			return err
		}

		// check for /profiles/ in the url, since profiles with vanity urls
		// set will always use /id/ instead
		dir := path.Dir(u.Path)
		if r, _ := path.Match("/profiles", dir); r {
			// path.Base(u.Path) gets the steamID64 from the end of the url
			convertedIDs, _ := convertSteamID(skey, path.Base(u.Path))
			printSteamIDs(convertedIDs)
			return nil
		}

		// convert the vanity URL into steamID64
		vanityURL := path.Base(u.Path)
		steamID64, err := steamid.ConvertVanityURL(vkey, vanityURL)
		if err != nil {
			return err
		}

		convertedIDs, _ := convertSteamID(skey, steamID64)
		if err != nil {
			return err
		}

		printSteamIDs(convertedIDs)
		return nil
	},
}

func printSteamIDs(data response) {
	msg := fmt.Sprintf(`
|  steam3ID: %s
| steamID32: %s
| steamID64: %s
	`, data.Converted.Steam3, data.Converted.SteamID32, data.Converted.SteamID64)

	fmt.Print(msg)
}

func convertSteamID(skey, steamID64 string) (response, error) {

	apiURL, err := url.Parse("https://steamidapi.uk/convert.php")

	q := apiURL.Query()
	if err != nil {
		return response{}, err
	}

	// build our request
	q.Set("api", skey)
	q.Set("input", steamID64)
	q.Set("format", "json")
	apiURL.RawQuery = q.Encode()

	r, err := http.Get(apiURL.String())
	if err != nil {
		return response{}, err
	}
	response := response{}

	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &response)

	return response, nil
}

func init() {
	rootCmd.AddCommand(idCmd)

	idCmd.Flags().StringVarP(&skey, "skey", "s", "", "steamid.uk API key (overrides config file, optional)")
	if err := viper.BindPFlag("skey", idCmd.Flags().Lookup("skey")); err != nil {
		log.Fatal("unable to bind flag: ", err)
	}
	idCmd.Flags().StringVarP(&vkey, "vkey", "v", "", "valve API key (overrides config file, optional)")
	if err := viper.BindPFlag("vkey", idCmd.Flags().Lookup("vkey")); err != nil {
		log.Fatal("unable to bind flag: ", err)
	}
	idCmd.Flags().StringVarP(&steamurl, "url", "u", "", "link to user's profile (required)")
	idCmd.MarkFlagRequired("url")
}
