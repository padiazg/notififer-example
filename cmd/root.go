/*
Copyright © 2024 Patricio Diaz <padiazg@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/padiazg/notifier-example/config/settings"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	s       = &settings.Settings{}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "notifier-example",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	listenCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notifier-example.yaml)")
	rootCmd.PersistentFlags().BoolP("show-config", "c", false, "Show config values")
	rootCmd.PersistentFlags().BoolP("show-key-values", "k", false, "Show key-value pairs")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	if err := s.Read(cfgFile); err != nil {
		panic(err)
	}
}
