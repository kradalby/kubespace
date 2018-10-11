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

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Output configuration for GitLab CI",
	Long: `Output configuration for GitLab CI.
Describes fields and and content to add to GitLab CI.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubespace.NewClient(kubeconf)
		if err != nil {
			panic(err)
		}

		endpoint := client.GetEndpoint()

		certificate, err := client.GetCertificate(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		token, err := client.GetToken(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		fmt.Println("Add the following to your Gitlab under 'Kubernetes cluster details'")
		fmt.Println()
		fmt.Println("Kubernetes cluster name:")
		fmt.Println("Descriptive name of your choice (e.g. gitlab-cluster)")
		fmt.Println()
		fmt.Println("API URL:")
		fmt.Println(endpoint)
		fmt.Println()
		fmt.Println("CA Certifcate:")
		fmt.Println(certificate)
		fmt.Println()
		fmt.Println("Token:")
		fmt.Println(token)
		fmt.Println()
		fmt.Println("Project namespace (optional, unique):")
		fmt.Println(namespace)

	},
}

func init() {
	rootCmd.AddCommand(gitlabCmd)

	gitlabCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (required)")
	gitlabCmd.MarkFlagRequired("namespace")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitlabCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitlabCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
