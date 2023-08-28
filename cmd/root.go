/*
Copyright © 2023 Arush Salil

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
	"fmt"
	"os"

	"github.com/arush-sal/gese/pkg/constant"
	"github.com/arush-sal/gese/pkg/encrypter"
	"github.com/arush-sal/gese/pkg/runner"
	"github.com/arush-sal/gese/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var env, token, repo, owner, secret, pubkey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:        "gese",
	SuggestFor: []string{"gasa", "gesr", "gwsw", "gede"},
	Short:      "GitHub Environment Secret Encrypter",
	Long: `GeSe pronounced as Gee-Sea
Get libsodium encrypted values that can be used for creating GitHub environment secrets`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var esecret string
		var err error

		if util.IsEmptyString(repo, owner) && util.IsEmptyString(pubkey) {
			return cmd.Help()
		}

		if !util.IsEmptyString(repo, owner) {
			if util.IsEmptyString(token) {
				token = os.Getenv(constant.TOKEN_NAME)
			}
			if util.IsEmptyString(token) {
				if term.IsTerminal(int(os.Stdin.Fd())) {
					fmt.Print("Enter GitHub Token: ")
					t, err := term.ReadPassword(0)
					if err != nil {
						return err
					}
					token = string(t)
					fmt.Println("")
				}
			}

			esecret, err = runner.GetEncryptedValue(token, repo, owner, secret, env)
			if err != nil {
				return err
			}
		}

		if !util.IsEmptyString(pubkey) {
			esecret, err = encrypter.Encrypt(pubkey, secret)
			if err != nil {
				return err
			}
		}

		if !util.IsEmptyString(esecret) {
			fmt.Println(esecret)
		}

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "Name of the the GitHub environment")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Name of the the GitHub repository")
	rootCmd.PersistentFlags().StringVarP(&owner, "owner", "o", "", "Name of the the GitHub repository owner")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitHub token to be used for all the communication\n\"GITHUB_TOKEN\" environment variable can be used instead")
	rootCmd.PersistentFlags().StringVarP(&token, "secret", "s", "", "GitHub environment secret value")
	rootCmd.PersistentFlags().StringVarP(&pubkey, "pubkey", "p", "", "Public Key of the GitHub environment\n(mutually exclusive with token and repo info)")
}
