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

package clonecmd

import (
	"context"
	"fmt"

	//docs "github.com/pkgserver-dev/pkgserver/internal/docs/generated/initdocs"

	pkgv1alpha1 "github.com/pkgserver-dev/pkgserver/apis/pkg/v1alpha1"
	"github.com/pkgserver-dev/pkgserver/apis/pkgrevid"
	"github.com/pkgserver-dev/pkgserver/pkg/client"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, version string, cfg *genericclioptions.ConfigFlags) *Runner {
	r := &Runner{}
	cmd := &cobra.Command{
		Use:  "clone TGT-PKGREV[<Target>.<REPO>.<REALM>.<PACKAGE>.<WORKSPACE>] SRC-PKGREV[<Target>.<REPO>.<REALM>.<PACKAGE>.<WORKSPACE>] [flags]",
		Args: cobra.ExactArgs(2),
		//Short:   docs.InitShort,
		//Long:    docs.InitShort + "\n" + docs.InitLong,
		//Example: docs.InitExamples,
		PreRunE: r.preRunE,
		RunE:    r.runE,
	}

	r.Command = cmd
	r.cfg = cfg

	return r
}

func NewCommand(ctx context.Context, version string, kubeflags *genericclioptions.ConfigFlags) *cobra.Command {
	return NewRunner(ctx, version, kubeflags).Command
}

type Runner struct {
	Command *cobra.Command
	cfg     *genericclioptions.ConfigFlags
	client  client.Client
}

func (r *Runner) preRunE(_ *cobra.Command, _ []string) error {
	client, err := client.CreateClientWithFlags(r.cfg)
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	ctx := c.Context()
	//log := log.FromContext(ctx)
	//log.Info("create packagerevision", "src", args[0], "dst", args[1])

	namespace := "default"
	if r.cfg.Namespace != nil && *r.cfg.Namespace != "" {
		namespace = *r.cfg.Namespace
	}

	key := types.NamespacedName{
		Namespace: namespace,
		Name:      args[1],
	}
	srcPkgrev := &pkgv1alpha1.PackageRevision{}
	if err := r.client.Get(ctx, key, srcPkgrev); err != nil {
		return fmt.Errorf("cannot get src pkgRev from arg[1], err: %s", err.Error())
	}

	pkgRevName := args[0]
	dstPkgRevID, err := pkgrevid.ParsePkgRev2PkgRevID(pkgRevName)
	if err != nil {
		return err
	}

	pkgRev := pkgv1alpha1.BuildPackageRevision(
		metav1.ObjectMeta{
			Name:      pkgRevName,
			Namespace: namespace,
		},
		pkgv1alpha1.PackageRevisionSpec{
			PackageRevID: *dstPkgRevID,
			Lifecycle:    pkgv1alpha1.PackageRevisionLifecycleDraft,
			Upstream: &pkgrevid.Upstream{
				Repository: srcPkgrev.Spec.PackageRevID.Repository,
				Realm:      srcPkgrev.Spec.PackageRevID.Realm,
				Package:    srcPkgrev.Spec.PackageRevID.Package,
				Revision:   srcPkgrev.Spec.PackageRevID.Revision,
			},
			Tasks: []pkgv1alpha1.Task{
				{
					Type: pkgv1alpha1.TaskTypeClone,
				},
			},
		},
		pkgv1alpha1.PackageRevisionStatus{},
	)

	if err := r.client.Create(ctx, pkgRev); err != nil {
		return err
	}
	fmt.Println(pkgRev.Name)
	return nil
}
