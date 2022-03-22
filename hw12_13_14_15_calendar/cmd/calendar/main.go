/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/cmd"
	"github.com/alexei38/otus_hw/hw12_13_14_15_calendar/pkg/cmd/cli/calendar"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Running the calendar application",

	Version: cmd.GetVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		calendar.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.Flags().String(
		"config",
		"",
		"config file (default is $HOME/.calendar.yaml, /etc/calendar.yaml)",
	)
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
}

func main() {
	Execute()
}
