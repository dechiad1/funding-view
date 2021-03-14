/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"company-funding/parser"
	"company-funding/util"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var (
	year  int
	fetch bool

	// batchLoadCmd represents the batchLoad command
	batchLoadCmd = &cobra.Command{
		Use:   "batchLoad",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("batchLoad called")
			if year == 0 {
				fmt.Println("A year to run the batch load for is required")
				os.Exit(1)
			}

			if fetch {
				//request.BatchFetchWebpages(year)
			}

			dir := fmt.Sprintf("webpages/fundingletter/%d/", year)
			fmt.Printf("fetching from directory %s\n", dir)
			files := util.GetFilenamesFromDirectory(dir)
			for _, f := range files {
				fmt.Printf("Attempting to parse file %s\n", dir+f)
				b, err := util.GetLocalHtml(dir + f)
				if err != nil {
					panic(err)
				}

				doc, _ := html.Parse(strings.NewReader(string(b)))
				p := parser.FlParser{}
				p.Parse(doc)
			}
		},
	}
)

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batchLoadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// batchLoadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	batchLoadCmd.Flags().BoolP("fetch", "f", false, "fetch as well as the load the data for a year")
	batchLoadCmd.Flags().IntVar(&year, "year", 0, "the year to run the batch process for")
	rootCmd.AddCommand(batchLoadCmd)
}
