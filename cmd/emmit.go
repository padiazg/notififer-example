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
	"github.com/padiazg/notifier-example/emitter"
	"github.com/spf13/cobra"
)

// emmitCmd represents the emmit command
var emmitCmd = &cobra.Command{
	Use:   "emmit",
	Short: "Sends a few messages",
	Long:  `Notifies a few payloads to the different channels set in the configuration file, then exits immediately`,
	Run: func(cmd *cobra.Command, args []string) {
		s.ParseServeArgsAndFlags(cmd, args)
		s.Show()
		emitter.Run(s)
	},
}

func init() {
	rootCmd.AddCommand(emmitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// emmitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// emmitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
