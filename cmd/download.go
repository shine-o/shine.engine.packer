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
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download .gdp files which contain patch data",
	Run: handlers.Download,

}

func init() {
	rootCmd.AddCommand(downloadCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	downloadCmd.PersistentFlags().String("patch-hive", "http://patch.cdn.gamigo.com/fo/es/PatchHive.txt", "list of .gdp files, e.g: PatchHive.txt, http://patch.cdn.gamigo.com/fous/gdp/PatchHive.txt")
	downloadCmd.PersistentFlags().String("destination", "./downloaded", "Path on system where downloaded .gdp files should be persisted")
	downloadCmd.PersistentFlags().Bool("overwrite", false, "Download and replace if the file already exists")
}
