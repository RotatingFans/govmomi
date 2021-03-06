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

package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/RotatingFans/govmomi/object"
	"golang.org/x/net/context"
)

type FolderFlag struct {
	common

	*DatacenterFlag

	name   string
	folder *object.Folder
}

var folderFlagKey = flagKey("folder")

func NewFolderFlag(ctx context.Context) (*FolderFlag, context.Context) {
	if v := ctx.Value(folderFlagKey); v != nil {
		return v.(*FolderFlag), ctx
	}

	v := &FolderFlag{}
	v.DatacenterFlag, ctx = NewDatacenterFlag(ctx)
	ctx = context.WithValue(ctx, folderFlagKey, v)
	return v, ctx
}

func (flag *FolderFlag) Register(ctx context.Context, f *flag.FlagSet) {
	flag.RegisterOnce(func() {
		flag.DatacenterFlag.Register(ctx, f)

		env := "GOVC_FOLDER"
		value := os.Getenv(env)
		usage := fmt.Sprintf("Folder [%s]", env)
		f.StringVar(&flag.name, "folder", value, usage)
	})
}

func (flag *FolderFlag) Process(ctx context.Context) error {
	return flag.ProcessOnce(func() error {
		if err := flag.DatacenterFlag.Process(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (flag *FolderFlag) Folder() (*object.Folder, error) {
	if flag.folder != nil {
		return flag.folder, nil
	}

	finder, err := flag.Finder()
	if err != nil {
		return nil, err
	}

	if flag.folder, err = finder.FolderOrDefault(context.TODO(), flag.name); err != nil {
		return nil, err
	}

	return flag.folder, nil
}

func (flag *FolderFlag) FolderOrRoot() (*object.Folder, error) {
	if flag.name == "" {
		client, err := flag.Client()
		if err != nil {
			return nil, err
		}

		flag.folder = object.NewRootFolder(client)
		return flag.folder, nil
	}

	return flag.Folder()
}
