// Copyright (c) 2019, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"os"

	sendmail "github.com/retailnext/sentry-sendmail"
)

func main() {
	sendmail.ParseOptions()
	err := sendmail.SentryConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	headers, body := sendmail.ReadData(bufio.NewReader(os.Stdin))
	message, err := sendmail.BuildMessage(headers, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sendmail.SentrySend(message, headers)
	if err != nil {
		fmt.Println(err)
		return
	}
}
