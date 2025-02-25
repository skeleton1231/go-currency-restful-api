// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	apiclientv1 "github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam/apiserver/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/skeleton1231/go-iam-ecommerce-microservice/internal/iamctl/cmd/util"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/internal/iamctl/util/templates"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/pkg/cli/genericclioptions"
)

const (
	createUsageStr = "create SECRET_NAME"
)

// CreateOptions is an options struct to support create subcommands.
type CreateOptions struct {
	Description string
	Expires     int64

	Secret *v1.Secret

	Client apiclientv1.APIV1Interface

	genericclioptions.IOStreams
}

var (
	createLong = templates.LongDesc(`Create secret resource.

This will generate secretID and secretKey which can be used to sign JWT token.`)

	createExample = templates.Examples(`
		# Create secret which will expired after 2 hours
		iamctl secret create foo

		# Create secret with a specified expire time and description
		iamctl secret create foo --expires=1988121600 --description="secret for iam"`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nSECRET_NAME is required arguments for the create command",
		createUsageStr,
	)
)

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions(ioStreams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		Expires:   time.Now().Add(144 * time.Hour).Unix(),
		IOStreams: ioStreams,
	}
}

// NewCmdCreate returns new initialized instance of create sub command.
func NewCmdCreate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Create secret resource",
		TraverseChildren:      true,
		Long:                  createLong,
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Description, "description", o.Description, "The descriptin of the secret.")
	cmd.Flags().Int64Var(&o.Expires, "expires", o.Expires, "The expire time of the secret.")

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	o.Secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Expires:     o.Expires,
		Description: o.Description,
	}

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = apiclientv1.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if errs := o.Secret.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

// Run executes a create subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	secret, err := o.Client.Secrets().Create(context.TODO(), o.Secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "secret/%s created\n", secret.Name)

	return nil
}
