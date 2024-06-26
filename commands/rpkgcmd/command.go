/*
Copyright 2024 Nokia.

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

package rpkgcmd

import (
	"context"

	//docs "github.com/pkgserver-dev/pkgserver/internal/docs/generated/initdocs"

	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/clonecmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/createcmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/deletecmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/getcmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/pullcmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/pushcmd"
	"github.com/pkgserver-dev/pkgctl/commands/rpkgcmd/updatestatuscmd"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewRunner returns a command runner.
func GetCommand(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use: "rpkg",
		//Short:   docs.InitShort,
		//Long:    docs.InitShort + "\n" + docs.InitLong,
		//Example: docs.InitExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cmd.Flags().GetBool("help")
			if err != nil {
				return err
			}
			if h {
				return cmd.Help()
			}
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		//approvecmd.NewCommand(ctx, version, kubeflags),
		//clonecmd.NewCommand(ctx, version, kubeflags),
		getcmd.NewCommand(ctx, version, kubeflags),
		createcmd.NewCommand(ctx, version, kubeflags),
		clonecmd.NewCommand(ctx, version, kubeflags),
		deletecmd.NewCommand(ctx, version, kubeflags),
		updatestatuscmd.NewCommand(ctx, version, kubeflags),
		pullcmd.NewCommand(ctx, version, kubeflags),
		pushcmd.NewCommand(ctx, version, kubeflags),
		////proposecmd.NewCommand(ctx, version, kubeflags),
		////proposedeletecmd.NewCommand(ctx, version, kubeflags),
		//pushcmd.NewCommand(ctx, version, kubeflags),
	)
	return cmd
}
