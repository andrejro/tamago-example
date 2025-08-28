// Copyright (c) The TamaGo Authors. All Rights Reserved.
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

//go:build mx6ullevk || usbarmory || cloud_hypervisor || microvm || firecracker

package cmd

import (
	"fmt"
	"net"
	"regexp"

	"github.com/usbarmory/tamago-example/shell"
)

func init() {
	shell.Add(shell.Cmd{
		Name:    "dns",
		Args:    1,
		Pattern: regexp.MustCompile(`^dns (.*)`),
		Syntax:  "<host>",
		Help:    "resolve domain",
		Fn:      dnsCmd,
	})
}

func dnsCmd(_ *shell.Interface, arg []string) (res string, err error) {
	cname, err := net.LookupHost(arg[0])

	if err != nil {
		return "", fmt.Errorf("query error: %v", err)
	}

	return fmt.Sprintf("%+v", cname), nil
}
