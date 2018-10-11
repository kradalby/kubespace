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

var repo string

// droneCmd represents the drone command
var droneCmd = &cobra.Command{
	Use:   "drone",
	Short: "Output commands for adding secrets to drone CI",
	Long: `Output commands for adding secrets to drone CI
	
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubespace.NewClient(kubeconf)
		if err != nil {
			panic(err)
		}

		endpoint := client.GetEndpoint()

		certificate, err := client.GetCertificateB64(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		token, err := client.GetToken(namespace)
		if err != nil {
			log.Printf("[Error] %s", err)
		}

		fmt.Printf("drone secret add %s --image=quay.io/honestbee/drone-kubernetes --name=kubernetes_server --value=%s\n", repo, endpoint)
		fmt.Printf("drone secret add %s --image=quay.io/honestbee/drone-kubernetes --name=kubernetes_cert --value=%s\n", repo, certificate)
		fmt.Printf("drone secret add %s --image=quay.io/honestbee/drone-kubernetes --name=kubernetes_token --value=%s\n", repo, token)
	},
}

func init() {
	rootCmd.AddCommand(droneCmd)

	droneCmd.Flags().StringVarP(&repo, "repo", "r", "", "drone repository to add secrets (required)")
	droneCmd.MarkFlagRequired("repo")

	droneCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (required)")
	droneCmd.MarkFlagRequired("namespace")
}
