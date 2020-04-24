/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/shine-o/shine.engine.packer/handlers"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract .gdp files into separate folders",
	Run: handlers.Extract,

}

func init() {
	rootCmd.AddCommand(extractCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	extractCmd.PersistentFlags().String("source", "./downloaded", "path where raw .gdp files are stored (default is ./downloaded)")
	extractCmd.PersistentFlags().String("destination", "./extracted", "path where extracted data should be stored (default is ./extracted)")
	//extractCmd.PersistentFlags().Bool("accumulative", false, "Store files in a single location, similar to building a client but without the base (installer files) (default is false")
	extractCmd.PersistentFlags().Bool("server-files", false, "Store potential server files(if found) in a separate location (default is false")
	extractCmd.PersistentFlags().String("server-files-path", "./potential-server-files", "Location where to store potential server files (if found) in a separate location (default is false")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
