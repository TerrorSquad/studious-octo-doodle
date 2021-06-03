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
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type Product struct {
	sku    string
	images map[string]string
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bicsv",
	Short: "Generates a CSV files that contains data from bulk image import",
	Long: `The command accepts exactly one argument, path to the images directory.

It will generate a CSV string and output it to STDOUT when called.
CSV will have the following format:
sku, base_image, small_image, thumbnail_image, rollover_image

Example: bicsv ./product_images
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]

		files, err := ioutil.ReadDir(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		re := regexp.MustCompile(`(?P<sku>\d+)_(?P<suffix>\d)(?P<extension>\.jpg|png)`)
		var products = map[string]Product{}
		for _, file := range files {
			var match = re.FindAllStringSubmatch(file.Name(), -1)
			if len(match) > 0 {
				sku := match[0][1]
				suffix := match[0][2]
				if products[sku].sku == "" {
					products[sku] = updateImagesMapForProduct(Product{
						sku:    sku,
						images: map[string]string{},
					}, suffix, file.Name())
				} else {
					if sku == products[sku].sku {
						updateImagesMapForProduct(products[sku], suffix, file.Name())
					}
				}
			}
		}
		headers := []string{"sku", "base_image", "small_image", "thumbnail_image", "rollover_image"}
		var rows [][]string
		for _, product := range products {
			rows = append(rows, []string{
				product.sku,
				product.images["baseImage"],
				product.images["smallImage"],
				product.images["thumbnailImage"],
				product.images["rolloverImage"]})
		}
		byteData, err := WriteAll(append([][]string{headers}, rows...))
		fmt.Print(string(byteData))
	},
}

func updateImagesMapForProduct(product Product, suffix string, fileName string) Product {
	if suffix == "1" {
		product.images["baseImage"] = fileName
		product.images["smallImage"] = fileName
		product.images["thumbnailImage"] = fileName
	} else if suffix == "2" {
		product.images["rolloverImage"] = fileName
	}
	return product
}

func WriteAll(records [][]string) ([]byte, error) {
	if records == nil || len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	err := csvWriter.WriteAll(records)
	if err != nil {
		return nil, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bicsv.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bicsv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bicsv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
