// Copyright © 2018 Skyscrapers <hello@skyscrapers.eu`
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joshdk/ykmango"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var region string
var profile string
var serialNumber string
var slotName string
var clientListLocation string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gasy",
	Short: "Go-AWS-STS-YubiKey CLI tool",
	Long: `A CLI tool to generate STS keys and URLs using Yubikey OTP.

Please see the README for documentation: https://github.com/skyscrapers/gasy`,
	Args: cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		account := getAccount(args[0])

		// Generate an OATH code using the given slot name.
		// You may need to touch your YubiKey device if the
		// slot is configured to require touch.
		boldGreen := color.New(color.FgGreen, color.Bold)
		boldGreen.Println("Please touch your YubiKey...")

		code, err := ykman.Generate(slotName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		region := viper.Get("aws.region")
		boldGreen.Println("requesting credentials for " + account.Name)
		// Request a token from STS using the code
		login(region.(string), code, serialNumber, account)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gasy.toml)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "eu-west-1", "region to use with AWS")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "which AWS profile to use to perform the login (default is default)")
	rootCmd.PersistentFlags().StringVarP(&serialNumber, "serialnumber", "s", "", "serial number of your AWS MFA device")
	rootCmd.PersistentFlags().StringVarP(&slotName, "slotname", "S", "", "Name of your YubiKey ath slot")
	rootCmd.PersistentFlags().StringVarP(&clientListLocation, "client-list-location", "c", "", "Path to the json client list")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gasy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gasy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// populate defaults
	viper.SetDefault("aws.region", "eu-west-1")

	// if the profile is not set by a flag, use the one in the config file
	if profile == "" {
		viper.SetDefault("aws.profile", "default")
		profile = viper.Get("aws.profile").(string)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	//check if we have our required variables
	if serialNumber == "" {
		if viper.Get("aws.mfaSerial") == nil {
			fmt.Println("No mfa serialnumber configured")
			os.Exit(1)
		} else {
			serialNumber = viper.Get("aws.mfaSerial").(string)
		}
	}

	if slotName == "" {
		if viper.Get("yubikey.slotName") == nil {
			fmt.Println("No YubiKey slot name configured")
			os.Exit(1)
		} else {
			slotName = viper.Get("yubikey.slotName").(string)
		}
	}
	if clientListLocation == "" {
		if viper.Get("aws.clientListLocation") == nil {
			fmt.Println("No client list location configured")
			os.Exit(1)
		} else {
			clientListLocation = viper.Get("aws.clientListLocation").(string)
		}
	}

}
