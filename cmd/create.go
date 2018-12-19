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
	"fmt"
	"log"

	kubespace "github.com/kradalby/kubespace/kubespace"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a namespace with a restricted account",
	Long: `Create a namespace with a restricted account.
	
Creates:
* Namespace
* Service account
* Role with namespace permissions
* Binds service account to role
* Outputs kubeconf for user`,
	Run: func(cmd *cobra.Command, args []string) {
		// namespace := cmd.Flag("namespace").Value.String()
		// kubeconf := cmd.Flag("kubeconf").Value.String()
		// skipns := cmd.Flag("skipns").Value
		client, err := kubespace.NewClient(kubeconf)
		if err != nil {
			panic(err)
		}

		err = client.CreateIfNotExistClusterRole()
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		if !skipns {
			err = client.CreateNamespace(namespace)
			if err != nil {
				log.Printf("[Error] %s", err)
			}
		}

		err = client.CreateServiceAccount(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		err = client.CreateIfNotExistServiceAccountClusterRoleBinding(namespace)
		if err != nil {
			// log.Printf("[Error] %s", err)
			log.Println(err)
		}

		err = client.CreateServiceAccountRoleBinding(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		err = client.CreateRole(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		config, err := client.CreateConfiguration(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}
		fmt.Println()
		fmt.Println("Configuration for namespace:")
		fmt.Println(config)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (required)")
	createCmd.MarkFlagRequired("namespace")
	createCmd.Flags().BoolVarP(&skipns, "skip-namespace", "s", false, "Skip namespace creation")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
