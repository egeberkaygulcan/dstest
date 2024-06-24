package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var DEFAULT_COMMAND = "run"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dstest",
	Short: "Testing distributed systems, one Heisenbug at a time",
	Long: `dstest is a concurrency testing tool for distributed systems.
It is designed to test distributed systems with different
schedulers and network configurations without the need to
instrument the system under test.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// print hello world
		cmd.Println("Hello, World!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dstest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func subCommands() (commandNames []string) {
	for _, command := range rootCmd.Commands() {
		commandNames = append(commandNames, append(command.Aliases, command.Name())...)
	}
	return
}

func setDefaultCommandIfNonePresent() {
	if len(os.Args) > 1 {
		potentialCommand := os.Args[1]
		for _, command := range subCommands() {
			if command == potentialCommand {
				return
			}
		}
		os.Args = append([]string{os.Args[0], DEFAULT_COMMAND}, os.Args[1:]...)
	}

}
