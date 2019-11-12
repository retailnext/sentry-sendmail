// Copyright (c) 2019, RetailNext, Inc.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
// All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goreleaser/nfpm"
	_ "github.com/goreleaser/nfpm/deb"
)

var (
	// Version specifies Semantic versioning increment (MAJOR.MINOR.PATCH).
	Version = "0.0.0"
)

var yamlConfig = `# sentry-sendmail config file
name: "sentry-sendmail"
arch: "amd64"
platform: "linux"
version: "` + Version + `"
section: "default"
priority: "extra"
replaces:
- sendmail
maintainer: "Ivan Daunis <ivan.daunis@retailnext.net>"
description: |
  Sendmail interface for sentry.
vendor: "RetailNext Inc."
homepage: "http://github.com/retailnext/sentry-sendmail"
license: "BSD"
files:
  ./dist/sendmail: "/usr/sbin/sendmail"
config_files:
  ./config/sentry-sendmail.conf: "/etc/sentry-sendmail.conf"
`

func buildPackage(yaml, format, target string) error {
	config, err := nfpm.Parse(strings.NewReader(yaml))
	if err != nil {
		return err
	}

	info, err := config.Get(format)
	if err != nil {
		return err
	}

	info = nfpm.WithDefaults(info)

	if err = nfpm.Validate(info); err != nil {
		return err
	}

	pkg, err := nfpm.Get(format)
	if err != nil {
		return err
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}

	return pkg.Package(info, f)
}

func main() {
	err := buildPackage(yamlConfig, "deb", "dist/sentry-sendmail-v"+Version+".deb")
	if err != nil {
		fmt.Println(err)
	}
}
