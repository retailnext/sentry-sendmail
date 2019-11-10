// Copyright (c) 2019, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.

package main

import (
	"github.com/xor-gate/debpkg"
)

var (
	// Version specifies Semantic versioning increment (MAJOR.MINOR.PATCH).
	Version = "0.0.0"
)

var postInst = `#!/bin/sh
chown root:root /usr/sbin/sendmail
chown root:root /etc/sentry-sendmail.conf
`

func main() {
	deb := debpkg.New()

	deb.SetName("sentry-sendmail")
	deb.SetVersion(Version)
	deb.SetArchitecture("amd64")
	deb.SetMaintainer("Ivan Daunis")
	deb.SetMaintainerEmail("ivan.daunis@retailnext.net")
	deb.SetHomepage("http://github.com/retailnext")
	deb.SetShortDescription("Sendmail interface for sentry")
	deb.SetDescription("Sendmail interface for sentry\n")

	deb.AddControlExtraString("postinst", postInst)
	deb.AddFile("config/sentry-sendmail.conf", "./etc/sentry-sendmail.conf")
	deb.AddFile("dist/sendmail", "./usr/sbin/sendmail")

	deb.Write("dist/sentry-sendmail-v" + Version + ".deb")
	deb.Close()
}
