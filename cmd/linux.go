// Copyright (c) The TamaGo Authors. All Rights Reserved.
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

//go:build amd64

package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/u-root/u-root/pkg/boot/bzimage"

	"github.com/usbarmory/armory-boot/exec"
	"github.com/usbarmory/tamago-example/shell"
	"github.com/usbarmory/tamago/dma"
)

const (
	memoryStart = 0x80000000
	memorySize  = 0x10000000
	commandLine = "console=ttyS0,115200,8n1 mem=4G\x00"

	defaultPath = "bzImage"
)

var memoryMap = []bzimage.E820Entry{
	// should always be usable
	bzimage.E820Entry{
		Addr:    uint64(0x00000000),
		Size:    uint64(0x0009f000),
		MemType: bzimage.RAM,
	},
	// amd64.ramStart, microvm.ramSize
	bzimage.E820Entry{
		Addr:    0x10000000,
		Size:    0x40000000,
		MemType: bzimage.RAM,
	},
	bzimage.E820Entry{
		Addr:    memoryStart,
		Size:    memorySize,
		MemType: bzimage.RAM,
	},
}

func init() {
	shell.Add(shell.Cmd{
		Name:    "linux",
		Args:    1,
		Pattern: regexp.MustCompile(`^linux(.*)`),
		Syntax:  "(path)?",
		Help:    "boot Linux kernel bzImage",
		Fn:      linuxCmd,
	})
}

// For a full bare metal Go bootloader implementation see
// https://github.com/usbarmory/go-boot
func linuxCmd(_ *shell.Interface, arg []string) (res string, err error) {
	var bzImage []byte
	var mem *dma.Region

	path := strings.TrimSpace(arg[0])

	if len(path) == 0 {
		path = defaultPath
	}

	if bzImage, err = os.ReadFile(path); err != nil {
		return
	}

	if mem, err = dma.NewRegion(memoryStart, memorySize, false); err != nil {
		return
	}

	mem.Reserve(memorySize, 0)

	image := &exec.LinuxImage{
		Memory:  memoryMap,
		Region:  mem,
		Kernel:  bzImage,
		CmdLine: commandLine,
	}

	if err = image.Load(); err != nil {
		return "", fmt.Errorf("could not load kernel, %v", err)
	}

	log.Printf("starting kernel@%0.8x\n", image.Entry())

	return "", image.Boot(nil)
}
