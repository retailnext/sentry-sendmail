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

func formatError(err error) string {
	return fmt.Sprintf("%s: error: %s", appName, err.Error())
}

func errorWithMessage(err error, message string) {
	if journal.Enabled() {
		vars := map[string]string{}
		journal.Send(formatError(err), journal.PriErr, vars)
		journal.Send(message, journal.PriErr, vars)
	} else {
		bytes, _ := json.Marshal(map[string]string{"message": message})
		fmt.Fprintln(os.Stderr, formatError(err), string(bytes))
	}
}

func main() {
	appName = filepath.Base(os.Args[0])

	sendmail.ParseOptions()
	headers, body, rawData := sendmail.ReadData(bufio.NewReader(os.Stdin))

	err := sendmail.SentryConfig()
	if err != nil {
		errorWithMessage(err, rawData)
		return
	}

	message, err := sendmail.BuildMessage(headers, body)
	if err != nil {
		errorWithMessage(err, rawData)
		return
	}

	err = sendmail.SentrySend(message, headers)
	if err != nil {
		errorWithMessage(err, rawData)
		return
	}
}
