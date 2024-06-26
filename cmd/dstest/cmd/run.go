package cmd

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the test engine",
	Long:  `Run the test engine with the given configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		for key, value := range viper.GetViper().AllSettings() {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println("Starting dstest")

		// Read config
		cfg, err := config.Read()
		if err != nil {
			log.Fatal(err.Error())
		}

		// -----------------------------

		fmt.Println("Name: " + cfg.TestConfig.Name)

		te := new(engine.TestEngine)
		te.Init(cfg)

		te.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// path to configuration file
	runCmd.PersistentFlags().StringP("config", "c", "./config/config.yml", "Path to configuration file")
	err := viper.BindPFlag("config", runCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		log.Fatal(fmt.Errorf("error binding flag: %v", err))
		return
	}

	//runCmd.PersistentFlags().StringP("name", "n", "default", "Name of the test")
	//err = viper.BindPFlag("name", runCmd.PersistentFlags().Lookup("name"))
}
