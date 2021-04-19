package cmd

import (
	"fmt"
	"log"
	"net/url"
	"path"

	steamid "github.com/01dog/pcheck/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// used for flags
	skey     string
	vkey     string
	steamurl string
)

// idCmd represents the id command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "returns steamid information on a user",
	Long:  `long msg`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skey = viper.GetString("skey")
		vkey = viper.GetString("vkey")

		u, err := url.Parse(steamurl)
		if err != nil {
			fmt.Println("error parsing steam url: ", err)
			return err
		}

		// check if the user is using a vanity url
		// if not, there's no need to convert it first
		dir := path.Dir(u.Path)
		if r, _ := path.Match("/profiles", dir); r {
			fmt.Println("no need to convert")
			return nil
		}

		vanityURL := path.Base(u.Path)
		steamID64, err := steamid.ConvertVanityURL(vkey, vanityURL)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		fmt.Println(steamID64)
		return nil
	},
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
