/*
Copyright (c) 2014-2016 VMware, Inc. All Rights Reserved.

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

package datastore

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path"
	"text/tabwriter"

	"github.com/RotatingFans/govmomi/govc/cli"
	"github.com/RotatingFans/govmomi/govc/flags"
	"github.com/RotatingFans/govmomi/object"
	"github.com/RotatingFans/govmomi/units"
	"github.com/RotatingFans/govmomi/vim25/types"
	"golang.org/x/net/context"
)

type ls struct {
	*flags.DatastoreFlag
	*flags.OutputFlag

	long  bool
	slash bool
	all   bool
}

func init() {
	cli.Register("datastore.ls", &ls{})
}

func (cmd *ls) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.DatastoreFlag, ctx = flags.NewDatastoreFlag(ctx)
	cmd.DatastoreFlag.Register(ctx, f)

	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)

	f.BoolVar(&cmd.long, "l", false, "Long listing format")
	f.BoolVar(&cmd.slash, "p", false, "Write a slash (`/') after each filename if that file is a directory")
	f.BoolVar(&cmd.all, "a", false, "Include entries whose names begin with a dot (.)")
}

func (cmd *ls) Process(ctx context.Context) error {
	if err := cmd.DatastoreFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.OutputFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *ls) Usage() string {
	return "[FILE]..."
}

func (cmd *ls) Run(ctx context.Context, f *flag.FlagSet) error {
	ds, err := cmd.Datastore()
	if err != nil {
		return err
	}

	b, err := ds.Browser(context.TODO())
	if err != nil {
		return err
	}

	args := f.Args()
	if len(args) == 0 {
		args = []string{""}
	}

	result := &listOutput{
		rs:  make([]types.HostDatastoreBrowserSearchResults, 0),
		cmd: cmd,
	}

	for _, arg := range args {
		spec := types.HostDatastoreBrowserSearchSpec{
			MatchPattern: []string{"*"},
		}

		if cmd.long {
			spec.Details = &types.FileQueryFlags{
				FileType:     true,
				FileSize:     true,
				FileOwner:    types.NewBool(true), // TODO: omitempty is generated, but seems to be required
				Modification: true,
			}
		}

		for i := 0; ; i++ {
			r, err := cmd.ListPath(b, arg, spec)
			if err != nil {
				// Treat the argument as a match pattern if not found as directory
				if i == 0 && types.IsFileNotFound(err) {
					spec.MatchPattern[0] = path.Base(arg)
					arg = path.Dir(arg)
					continue
				}

				return err
			}

			// Treat an empty result against match pattern as file not found
			if i == 1 && len(r.File) == 0 {
				return fmt.Errorf("File %s/%s was not found", r.FolderPath, spec.MatchPattern[0])
			}

			result.add(r)
			break
		}
	}

	return cmd.WriteResult(result)
}

func (cmd *ls) ListPath(b *object.HostDatastoreBrowser, path string, spec types.HostDatastoreBrowserSearchSpec) (types.HostDatastoreBrowserSearchResults, error) {
	var res types.HostDatastoreBrowserSearchResults

	path, err := cmd.DatastorePath(path)
	if err != nil {
		return res, err
	}

	task, err := b.SearchDatastore(context.TODO(), path, &spec)
	if err != nil {
		return res, err
	}

	info, err := task.WaitForResult(context.TODO(), nil)
	if err != nil {
		return res, err
	}

	res = info.Result.(types.HostDatastoreBrowserSearchResults)
	return res, nil
}

type listOutput struct {
	rs  []types.HostDatastoreBrowserSearchResults
	cmd *ls
}

func (o *listOutput) add(r types.HostDatastoreBrowserSearchResults) {
	res := r
	res.File = nil

	for _, f := range r.File {
		if f.GetFileInfo().Path[0] == '.' && !o.cmd.all {
			continue
		}

		if o.cmd.slash {
			if d, ok := f.(*types.FolderFileInfo); ok {
				d.Path += "/"
			}
		}

		res.File = append(res.File, f)
	}

	o.rs = append(o.rs, res)
}

// hasMultiplePaths returns whether or not the slice of search results contains
// results from more than one folder path.
func (o *listOutput) hasMultiplePaths() bool {
	if len(o.rs) == 0 {
		return false
	}

	p := o.rs[0].FolderPath

	// Multiple paths if any entry is not equal to the first one.
	for _, e := range o.rs {
		if e.FolderPath != p {
			return true
		}
	}

	return false
}

func (o *listOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.rs)
}

func (o *listOutput) Write(w io.Writer) error {
	// Only include path header if we're dealing with more than one path.
	includeHeader := false
	if o.hasMultiplePaths() {
		includeHeader = true
	}

	tw := tabwriter.NewWriter(w, 3, 0, 2, ' ', 0)
	for i, r := range o.rs {
		if includeHeader {
			if i > 0 {
				fmt.Fprintf(tw, "\n")
			}
			fmt.Fprintf(tw, "%s:\n", r.FolderPath)
		}
		for _, file := range r.File {
			info := file.GetFileInfo()
			if o.cmd.long {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", units.ByteSize(info.FileSize), info.Modification.Format("Mon Jan 2 15:04:05 2006"), info.Path)
			} else {
				fmt.Fprintf(tw, "%s\n", info.Path)
			}
		}
	}
	tw.Flush()
	return nil
}
