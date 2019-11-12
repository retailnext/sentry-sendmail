// Copyright (c) 2019, RetailNext, Inc.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
// All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sendmail "github.com/retailnext/sentry-sendmail"
)

func main() {
	appName := filepath.Base(os.Args[0]) + ":"

	sendmail.ParseOptions()
	err := sendmail.SentryConfig()
	if err != nil {
		fmt.Println(appName, "error:", err)
		return
	}

	headers, body, rawData := sendmail.ReadData(bufio.NewReader(os.Stdin))
	message, err := sendmail.BuildMessage(headers, body)
	if err != nil {
		fmt.Println(appName, "error:", err)
		return
	}

	err = sendmail.SentrySend(message, headers)
	if err != nil {
		fmt.Println(appName, "error:", err)
		queueDir := "/var/spool/sentry-sendmail/queue"
		os.MkdirAll(queueDir, os.ModePerm)
		logFile, err := ioutil.TempFile(queueDir, "")
		if err != nil {
			fmt.Println(appName, "error: Can not save envelope:", err)
			return
		}
		fmt.Fprint(logFile, rawData)
		fmt.Println(appName, "info: Original envelope saved on:", queueDir)
		return
	}
}
