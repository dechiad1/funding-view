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
	"company-funding/repository"
	"company-funding/util"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var (
	file    string
	devMode bool

	devfileCmd = &cobra.Command{
		Use:   "devfile",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			// empty string is a test file
			content, err := util.GetLocalHtml(file)
			if err != nil {
				if file == "" {
					fmt.Println("tried to open a test.txt file in root directory, failed")
				} else {
					fmt.Printf("No such file %s exists\n", file)
				}

				panic(err)
			}

			doc, _ := html.Parse(strings.NewReader(string(content)))
			p := parser.FlParser{
				Dev:        devMode,
				CurrentDoc: file,
			}
			if !p.Dev {
				p.Db = repository.Connect()
			}
			p.Parse(doc)
		},
	}
)

func init() {
	devfileCmd.Flags().StringVar(&file, "file", "", "specify the file to get")
	devfileCmd.Flags().BoolVar(&devMode, "devMode", false, "Dev mode flag. Do not save to db")
	rootCmd.AddCommand(devfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
