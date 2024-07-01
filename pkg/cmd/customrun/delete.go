// Copyright © 2023 The Tekton Authors.
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

package customrun

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func deleteCommand(p cli.Params) *cobra.Command {
	//opts := &options.DeleteOptions{Resource: "CustomRun", ForceDelete: false, ParentResource: "CustomRun", DeleteAllNs: false}
	f := genericclioptions.NewPrintFlags("delete")
	eg := `Delete CustomRun with name 'foo' in namespace 'bar':

    tkn customrun delete foo -n bar

or

    tkn cr rm foo -n bar
`

	c := &cobra.Command{
		Use:               "delete",
		//Aliases:           []string{"rm"},
		Short:             "Delete CustomRuns having a specific name in a namespace",
		Example:           eg,
		Args:              cobra.MinimumNArgs(1), // Requires at least one argument (customrun-name)
		//SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			crNames := args
			s := &cli.Stream{
				In:  cmd.InOrStdin(),
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			err := deleteCustomRuns(s, p, crNames)
			if err != nil {
				return err
			}

			for _, crName := range crNames {
				fmt.Fprintf(s.Out, "CustomRun %s deleted successfully\n", crName)
			}
			return nil
		},
	}

	f.AddFlags(c)
	return c
}

func deleteCustomRuns(s *cli.Stream, p cli.Params, crNames []string) error {
	cs, err := p.Clients()
	if err != nil {
		return fmt.Errorf("failed to create tekton client: %w", err)
	}

	for _, crName := range crNames {
		err := deleteSingleCustomRun(cs, p.Namespace(), crName)
		if err != nil {
			fmt.Fprintf(s.Err, "failed to delete CustomRun %s: %v\n", crName, err)
			return err
		}
	}

	return nil
}

func deleteSingleCustomRun(cs *cli.Clients, namespace, crName string) error {
	err := cs.Dynamic.Resource(customrunGroupResource).Namespace(namespace).Delete(context.TODO(), crName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete CustomRun %s: %w", crName, err)
	}
	return nil
}
