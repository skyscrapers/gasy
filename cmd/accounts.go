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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Accounts struct holds an array
// of Account
type Accounts struct {
	Accounts []Account `json:"accounts"`
}

// Account struct holds account information
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SID         string `json:"sid,omitempty"`
	Description string `json:"description"`
}

var showID bool

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List all available accounts",
	Long: `List all accounts in the configuration file

You can use the number as a reference on login.`,
	Run: func(cmd *cobra.Command, args []string) {
		accounts := getAccountList()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Name", "Description", "ID"})

		for index, account := range accounts.Accounts {
			table.Append([]string{strconv.Itoa(index), account.Name, account.Description, account.ID})
		}

		table.Render()
	},
}

func getAccountList() Accounts {
	jsonFile, err := os.Open(clientListLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var accounts Accounts

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &accounts)

	return accounts
}

func getAccount(accountNumber string) Account {
	accounts := getAccountList()

	account := accounts.Accounts
	accountInt, _ := strconv.Atoi(accountNumber)
	//TODO: check & error if there is actually something there
	return account[accountInt]
}

func init() {
	rootCmd.AddCommand(accountsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.
	// accountsCmd.Flags().BoolP("showID", "i", false, "Show the account IDs")
}
