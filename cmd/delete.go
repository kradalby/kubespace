// Copyright © 2018 Kristoffer Dalby <kradalby@kradalby.no>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	kubespace "github.com/kradalby/kubespace/kubespace"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a namespace and the restricted account",
	Long:  `Delete a namespace and the restricted account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubespace.NewClient(kubeconf)
		if err != nil {
			panic(err)
		}

		err = client.DeleteServiceAccount(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		err = client.DeleteServiceAccountRoleBinding(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		err = client.DeleteRole(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}
		if !skipns {
			if askForConfirmation(fmt.Sprintf("You are currently trying to delete %s, this will delete ALL resources in the namespace, please confirm", namespace)) {
				err = client.DeleteNamespace(namespace)
				if err != nil {
					log.Printf("[Error] %s", err)
				}
			} else {
				fmt.Printf("Delete of namespace %s canceled\n", namespace)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (required)")
	deleteCmd.MarkFlagRequired("namespace")
	deleteCmd.Flags().BoolVarP(&skipns, "skip-namespace", "s", false, "Skip namespace delete")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
