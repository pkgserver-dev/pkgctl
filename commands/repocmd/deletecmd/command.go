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

package deletecmd

import (
	"context"

	//docs "github.com/pkgserver-dev/pkgserver/internal/docs/generated/initdocs"

	"github.com/fatih/color"
	"github.com/henderiw/logger/log"
	"github.com/pkgserver-dev/pkgserver/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *Runner {
	r := &Runner{}
	cmd := &cobra.Command{
		Use:  "delete NAME",
		Args: cobra.ExactArgs(2),
		//Short:   docs.InitShort,
		//Long:    docs.InitShort + "\n" + docs.InitLong,
		//Example: docs.InitExamples,
		PreRunE: r.preRunE,
		RunE:    r.runE,
	}

	r.Command = cmd
	r.kubeflags = kubeflags

	r.Command.Flags().StringVarP(
		&r.secret, "secret", "", "", "secret used for accessing the repository")
	r.Command.Flags().BoolVarP(
		&r.deployment, "deployment", "d", false, "tags the repository as a deployment repository. packages in a deployment repository are considered for deployment dependeing on their lifecycle status")
	r.Command.Flags().StringVarP(
		&r.directory, "directory", "", "", "the directory within the repository")
	return r
}

func NewCommand(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *cobra.Command {
	return NewRunner(ctx, version, kubeflags).Command
}

type Runner struct {
	Command   *cobra.Command
	kubeflags *genericclioptions.ConfigFlags
	client    client.Client
	local     bool
	// dynamic input
	secret     string
	deployment bool
	directory  string
}

func (r *Runner) preRunE(_ *cobra.Command, _ []string) error {
	client, err := client.CreateClientWithFlags(r.kubeflags)
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	ctx := c.Context()
	log := log.FromContext(ctx)
	log.Info("create repository", "local", r.local)
	if r.local {
		repoName := args[0]
		if !r.local {

			delete(viper.Get("repos").(map[string]interface{}), repoName)

			if err := viper.WriteConfig(); err != nil {
				color.Red("Error writing config file: %s", err.Error())
				return err
			}
		}
	}

	return nil
}
