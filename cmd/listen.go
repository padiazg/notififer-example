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
	"github.com/padiazg/notifier-example/listener"
	"github.com/spf13/cobra"
)

var (
	// listenCmd represents the listen command
	listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts the listeners",
		Long:  `Starts all the listeners set and configured in the configuration file. At least one listener must be set.`,
		Run: func(cmd *cobra.Command, args []string) {
			s.ParseServeArgsAndFlags(cmd, args)
			s.Show()
			listener.Run(s)
		},
	}
)

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
