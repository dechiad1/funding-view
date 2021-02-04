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
	"company-funding/datasource"
	"company-funding/util"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// devfileCmd represents the devfile command
var devfileCmd = &cobra.Command{
	Use:   "devfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := util.GetDevHtml()
		if err != nil {
			panic(err)
		}

		doc, _ := html.Parse(strings.NewReader(string(content)))
		node := util.GetNodeFromText(doc, "Americas")
		if node == nil {
			fmt.Println("no body!")
			os.Exit(1)
		}

		// get notable html.Node near data
		blockquote := util.GetParentOfType(node, "blockquote")

		// get nodes that contain relevant information
		sibs := util.GetSiblingsOfType(blockquote, "p")

		// each company has 7 datapoints. n-1 of the sibs array are valid, so remove the last
		input := sibs[0 : len(sibs)-1]

		if len(input)%7 != 0 {
			fmt.Printf("invalid input, wanted an input divisible by 7, got %d\n", len(input))
			return
		}

		count := len(input) / 7
		// for every 7, create a company
		for i := 0; i < count; i++ {
			start := i * 7
			end := (i + 1) * 7
			c := datasource.ParseCompany(input[start:end])
			fmt.Println(c.Name)
			// TODO: write to sqllite or something
		}
	},
}

func init() {
	rootCmd.AddCommand(devfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
