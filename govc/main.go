/*
Copyright (c) 2014-2015 VMware, Inc. All Rights Reserved.

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

package main

import (
	"os"

	"github.com/RotatingFans/govmomi/govc/cli"

	_ "github.com/RotatingFans/govmomi/govc/about"
	_ "github.com/RotatingFans/govmomi/govc/cluster"
	_ "github.com/RotatingFans/govmomi/govc/datacenter"
	_ "github.com/RotatingFans/govmomi/govc/datastore"
	_ "github.com/RotatingFans/govmomi/govc/device"
	_ "github.com/RotatingFans/govmomi/govc/device/cdrom"
	_ "github.com/RotatingFans/govmomi/govc/device/floppy"
	_ "github.com/RotatingFans/govmomi/govc/device/scsi"
	_ "github.com/RotatingFans/govmomi/govc/device/serial"
	_ "github.com/RotatingFans/govmomi/govc/dvs"
	_ "github.com/RotatingFans/govmomi/govc/dvs/portgroup"
	_ "github.com/RotatingFans/govmomi/govc/env"
	_ "github.com/RotatingFans/govmomi/govc/events"
	_ "github.com/RotatingFans/govmomi/govc/extension"
	_ "github.com/RotatingFans/govmomi/govc/fields"
	_ "github.com/RotatingFans/govmomi/govc/folder"
	_ "github.com/RotatingFans/govmomi/govc/host"
	_ "github.com/RotatingFans/govmomi/govc/host/account"
	_ "github.com/RotatingFans/govmomi/govc/host/autostart"
	_ "github.com/RotatingFans/govmomi/govc/host/esxcli"
	_ "github.com/RotatingFans/govmomi/govc/host/firewall"
	_ "github.com/RotatingFans/govmomi/govc/host/maintenance"
	_ "github.com/RotatingFans/govmomi/govc/host/option"
	_ "github.com/RotatingFans/govmomi/govc/host/portgroup"
	_ "github.com/RotatingFans/govmomi/govc/host/service"
	_ "github.com/RotatingFans/govmomi/govc/host/storage"
	_ "github.com/RotatingFans/govmomi/govc/host/vnic"
	_ "github.com/RotatingFans/govmomi/govc/host/vswitch"
	_ "github.com/RotatingFans/govmomi/govc/importx"
	_ "github.com/RotatingFans/govmomi/govc/license"
	_ "github.com/RotatingFans/govmomi/govc/logs"
	_ "github.com/RotatingFans/govmomi/govc/ls"
	_ "github.com/RotatingFans/govmomi/govc/permissions"
	_ "github.com/RotatingFans/govmomi/govc/pool"
	_ "github.com/RotatingFans/govmomi/govc/vapp"
	_ "github.com/RotatingFans/govmomi/govc/version"
	_ "github.com/RotatingFans/govmomi/govc/vm"
	_ "github.com/RotatingFans/govmomi/govc/vm/disk"
	_ "github.com/RotatingFans/govmomi/govc/vm/guest"
	_ "github.com/RotatingFans/govmomi/govc/vm/network"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
