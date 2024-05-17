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

package pushcmd

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	//docs "github.com/pkgserver-dev/pkgserver/internal/docs/generated/initdocs"

	"github.com/henderiw/store"
	"github.com/kform-dev/kform/pkg/pkgio"
	pkgv1alpha1 "github.com/pkgserver-dev/pkgserver/apis/pkg/v1alpha1"
	"github.com/pkgserver-dev/pkgserver/apis/pkgrevid"
	"github.com/pkgserver-dev/pkgserver/pkg/client"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *Runner {
	r := &Runner{}
	cmd := &cobra.Command{
		Use:  "push PKGREV[<Target>.<REPO>.<REALM>.<PACKAGE>.<WORKSPACE>] [DIR] [flags]",
		Args: cobra.MinimumNArgs(1),
		//Short:   docs.InitShort,
		//Long:    docs.InitShort + "\n" + docs.InitLong,
		//Example: docs.InitExamples,
		PreRunE: r.preRunE,
		RunE:    r.runE,
	}

	r.Command = cmd
	r.kubeflags = kubeflags

	return r
}

func NewCommand(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *cobra.Command {
	return NewRunner(ctx, version, kubeflags).Command
}

type Runner struct {
	Command   *cobra.Command
	kubeflags *genericclioptions.ConfigFlags
	client    client.Client
}

func (r *Runner) preRunE(_ *cobra.Command, _ []string) error {
	client, err := client.CreateClientWithFlags(r.kubeflags)
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func (r *Runner) runE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	//log := log.FromContext(ctx)
	//log.Info("create packagerevision", "src", args[0], "dst", args[1])

	namespace := "default"
	if r.kubeflags.Namespace != nil && *r.kubeflags.Namespace != "" {
		namespace = *r.kubeflags.Namespace
	}

	pkgRevName := args[0]
	if _, err := pkgrevid.ParsePkgRev2PkgRevID(pkgRevName); err != nil {
		return err
	}

	var resources map[string]string
	var err error
	if len(args) > 1 {
		resources, err = readFromDir(args[1])
	} else {
		resources, err = readFromReader(ctx, cmd.InOrStdin())
	}
	if err != nil {
		return err
	}

	key := types.NamespacedName{
		Namespace: namespace,
		Name:      pkgRevName,
	}
	prr := &pkgv1alpha1.PackageRevisionResources{}
	if err := r.client.Get(ctx, key, prr); err != nil {
		return err
	}

	prr.Spec.Resources = resources
	return r.client.Update(ctx, prr)
}

func readFromDir(dir string) (map[string]string, error) {
	resources := map[string]string{}
	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		resources[rel] = string(contents)
		return nil
	}); err != nil {
		return nil, err
	}
	return resources, nil
}

func readFromReader(ctx context.Context, in io.Reader) (map[string]string, error) {
	reader := pkgio.YAMLReader{
		Reader: in,
		Path:   "stdin",
	}
	data, err := reader.Read(ctx)
	if err != nil {
		return nil, err
	}
	resources := make(map[string]string, data.Len(ctx))
	data.List(ctx, func(ctx context.Context, k store.Key, r *yaml.RNode) {
		resources[k.Name] = r.MustString()
	})
	return resources, nil
}
