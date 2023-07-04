/*
Copyright © 2023 maxgio92 me@maxgio.me

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

package find

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/maxgio92/wfind/internal/output"
	"github.com/maxgio92/wfind/pkg/find"
)

type options struct {
	*find.Options
}

// NewCmd returns a new find command.
func NewCmd() *cobra.Command {
	o := &options{
		Options: &find.Options{},
	}

	cmd := &cobra.Command{
		Use:   "wfind URL",
		Short: "Find folders and files in web sites using HTTP or HTTPS",
		Args:  cobra.MinimumNArgs(1),
		RunE:  o.Run,
	}

	cmd.Flags().StringVarP(&o.FilenameRegexp, "name", "n", ".+", "Base of file name (the path with the leading directories removed) exact pattern.")
	cmd.Flags().StringVarP(&o.FileType, "type", "t", "", "The file type")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Enable verbosity to log all visited HTTP(s) files")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", true, "Whether to examine entries recursing into directories. Disable to behave like GNU find -maxdepth=0 option.")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewCmd()
	output.ExitOnErr(cmd.Execute())
}

func (o *options) validate() error {
	if err := o.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}

	return nil
}

func (o *options) Run(_ *cobra.Command, args []string) error {
	var seed string
	if len(args) > 0 {
		seed = args[0]
	}

	o.SeedURLs = append(o.SeedURLs, seed)

	if err := o.validate(); err != nil {
		return err
	}

	finder := find.NewFind(
		find.WithSeedURLs(o.SeedURLs),
		find.WithFilenameRegexp(o.FilenameRegexp),
		find.WithFileType(o.FileType),
		find.WithRecursive(o.Recursive),
		find.WithVerbosity(o.Verbose),
	)

	found, err := finder.Find()
	if err != nil {
		return errors.Wrap(err, "error finding the file")
	}

	for _, v := range found.URLs {
		output.Print(v)
	}

	return nil
}
