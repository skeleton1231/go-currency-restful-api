// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package version print the client and server version information.
package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/marmotedu/component-base/pkg/json"
	"github.com/marmotedu/component-base/pkg/version"
	restclient "github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/spf13/cobra"

	cmdutil "github.com/skeleton1231/go-iam-ecommerce-microservice/internal/iamctl/cmd/util"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/internal/iamctl/util/templates"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/pkg/cli/genericclioptions"
)

// Version is a struct for version information.
type Version struct {
	ClientVersion *version.Info `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion *version.Info `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
}

var versionExample = templates.Examples(`
		# Print the client and server versions for the current context
		iamctl version`)

// Options is a struct to support version command.
type Options struct {
	ClientOnly bool
	Short      bool
	Output     string

	client *restclient.RESTClient
	genericclioptions.IOStreams
}

// NewOptions returns initialized Options.
func NewOptions(ioStreams genericclioptions.IOStreams) *Options {
	return &Options{
		IOStreams: ioStreams,
	}
}

// NewCmdVersion returns a cobra command for fetching versions.
func NewCmdVersion(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(ioStreams)
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().BoolVar(
		&o.ClientOnly,
		"client",
		o.ClientOnly,
		"If true, shows client version only (no server required).",
	)
	cmd.Flags().BoolVar(&o.Short, "short", o.Short, "If true, print just the version number.")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "One of 'yaml' or 'json'.")

	return cmd
}

// Complete completes all the required options.
func (o *Options) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	var err error
	if o.ClientOnly {
		return nil
	}

	o.client, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate validates the provided options.
func (o *Options) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New(`--output must be 'yaml' or 'json'`)
	}

	return nil
}

// Run executes version command.
func (o *Options) Run() error {
	var (
		serverVersion *version.Info
		serverErr     error
		versionInfo   Version
	)

	clientVersion := version.Get()
	versionInfo.ClientVersion = &clientVersion

	if !o.ClientOnly && o.client != nil {
		// Always request fresh data from the server
		if err := o.client.Get().AbsPath("/version").Do(context.TODO()).Into(&serverVersion); err != nil {
			return err
		}
		versionInfo.ServerVersion = serverVersion
	}

	switch o.Output {
	case "":
		if o.Short {
			fmt.Fprintf(o.Out, "Client Version: %s\n", clientVersion.GitVersion)

			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", serverVersion.GitVersion)
			}
		} else {
			fmt.Fprintf(o.Out, "Client Version: %s\n", fmt.Sprintf("%#v", clientVersion))
			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", fmt.Sprintf("%#v", *serverVersion))
			}
		}
	case "yaml":
		marshaled, err := yaml.Marshal(&versionInfo)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.Out, string(marshaled))
	case "json":
		marshaled, err := json.MarshalIndent(&versionInfo, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(o.Out, string(marshaled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("VersionOptions were not validated: --output=%q should have been rejected", o.Output)
	}

	return serverErr
}
