package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var key string

// idCmd represents the id command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "returns steamid information on a user",
	Long: `long msg`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key = viper.GetString("key")
		return test(key)
	},
}

func test(key string) error {
	fmt.Println(key)
	return nil
}

func init() {
	rootCmd.AddCommand(idCmd)

	idCmd.Flags().StringVarP(&key, "key", "k", "", "steamid.uk API key (overrides config file, optional)")
	if err := viper.BindPFlag("key", idCmd.Flags().Lookup("key")); err != nil {
		log.Fatal("unable to bind flag: ", err)
	}
}
