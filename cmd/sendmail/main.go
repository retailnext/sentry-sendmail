// Copyright (c) 2019, RetailNext, Inc.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
// All rights reserved.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/journal"

	sendmail "github.com/retailnext/sentry-sendmail"
)

var (
	appName string
)

func wrap(err error) string {
	return fmt.Sprintf("%s: error: %s", appName, err.Error())
}

func main() {
	appName = filepath.Base(os.Args[0])

	sendmail.ParseOptions()
	err := sendmail.SentryConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, wrap(err))
		return
	}

	headers, body, rawData := sendmail.ReadData(bufio.NewReader(os.Stdin))
	message, err := sendmail.BuildMessage(headers, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, wrap(err))
		return
	}

	err = sendmail.SentrySend(message, headers)
	if err != nil {
		if journal.Enabled() {
			vars := map[string]string{}
			journal.Send(wrap(err), journal.PriErr, vars)
			journal.Send(rawData, journal.PriErr, vars)
		} else {
			bytes, _ := json.Marshal(map[string]string{"message": rawData})
			fmt.Fprintln(os.Stderr, wrap(err), string(bytes))
		}
		return
	}
}
