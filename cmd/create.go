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
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createFile string
	// createCmd represents the create command
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates an example configuration file",
		Long: `This commnand will create a .notifier-example.yaml file with some sample data. Adjust the values
according to your environment`,
		Run: func(cmd *cobra.Command, args []string) {
			fileName, _ := cmd.Flags().GetString("config")
			fmt.Printf("create: %s\n", fileName)
			s.SaveExample(fileName)
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().StringVar(&createFile, "config", ".notifier-example.yaml", "config file (default is $CWD/.notifier-example.yaml)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
