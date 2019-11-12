// Copyright (c) 2019, RetailNext, Inc.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
// All rights reserved.

package sendmail

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/getsentry/raven-go"
)

var (
	// GitCommit specifies the git commit sha, set by the compiler.
	GitCommit = ""
	// Version specifies Semantic versioning increment (MAJOR.MINOR.PATCH).
	Version = "0.0.0"
)

var (
	logPath    = "/var/log/sentry-sendmail.log"
	configPath = "/etc/sentry-sendmail.conf"
	recipients = ""
)

type Config struct {
	SentryDSN   string `toml:"DSN"`
	Environment string
}

func SentryConfig() error {
	var conf Config

	_, err := toml.DecodeFile(configPath, &conf)

	if os.Getenv("SENTRY_DSN") != "" {
		conf.SentryDSN = os.Getenv("SENTRY_DSN")
	}

	if os.Getenv("SENTRY_ENVIRONMENT") != "" {
		conf.Environment = os.Getenv("SENTRY_ENVIRONMENT")
	}

	// Error parsing the config file and we still don't have a DSN
	if _, isFileErr := err.(*os.PathError); conf.SentryDSN == "" && err != nil && !isFileErr {
		return fmt.Errorf("Can not read Sentry DSN from %s: %v", configPath, err)
	}

	if conf.SentryDSN == "" {
		return fmt.Errorf("Sentry DSN not set. Please set SENTRY_DSN environment variable or enter DSN in the config file: %s", configPath)
	}

	err = raven.SetDSN(conf.SentryDSN)
	if err != nil {
		return fmt.Errorf("Sentry DSN [%s] error: %v", conf.SentryDSN, err)
	}

	if conf.Environment != "" {
		raven.SetEnvironment(conf.Environment)
	}
	return nil
}

func getExtra(headers map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"ppid":    os.Getppid(),
		"headers": headers,
	}
}

func SentrySend(message string, headers map[string]string) error {
	strLevel := "Info"
	pkt := raven.NewPacketWithExtra(message, getExtra(headers))
	pkt.Level = raven.Severity(strLevel)
	eventID, ch := raven.Capture(pkt, nil)
	if eventID != "" {
		err := <-ch
		return err
	}
	return fmt.Errorf("Capture returned empty eventID")
}

func ReadData(reader *bufio.Reader) (map[string]string, string, string) {
	raw := ""
	body := ""
	headers := make(map[string]string)

	var logFile *os.File
	if opts.LogFile != "" {
		logFile, _ = os.OpenFile(opts.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintln(logFile, time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), "-----------")
		defer logFile.Close()
	}

	// Read headers first
	inHeaders := true
	prevKey := ""

	for {
		line, err := reader.ReadString('\n')
		if logFile != nil {
			fmt.Fprint(logFile, line)
		}
		raw += line
		if err != nil {
			body += line
			// If EOF here, line has not been processed and can be the body
			return headers, body, raw
		}
		if !opts.IgnoreDot && len(line) == 2 && line[0] == '.' {
			break
		}

		// Empty line indicates message body
		if len(line) == 1 {
			inHeaders = false
		}

		if inHeaders {
			index := strings.Index(line, ":")
			if index >= 0 || len(prevKey) > 0 {
				if index < 0 {
					headers[prevKey] += "\n" + strings.TrimSpace(line)
				} else {
					key := strings.ToLower(line[:index])
					value := strings.TrimSpace(line[index+1:])
					headers[key] = value
					prevKey = key
				}
			} else {
				inHeaders = false
			}
		}

		if !inHeaders {
			body += line
		}
	}

	return headers, body, raw
}

func BuildMessage(headers map[string]string, body string) (string, error) {
	message := headers["subject"] + "\n"

	if !opts.ExtractRecipients {
		headers["to"] = recipients
	}
	if len(opts.LegacyFrom) > 0 {
		headers["from"] = opts.LegacyFrom
	}
	if len(opts.SenderAddress) > 0 {
		headers["from"] = opts.SenderAddress
	}

	if headers["content-transfer-encoding"] == "quoted-printable" {
		decoded, _ := ioutil.ReadAll(quotedprintable.NewReader(strings.NewReader(string(body))))
		message += string(decoded)
	} else {
		message += string(body)
	}

	if len(headers["from"]) == 0 {
		return message, fmt.Errorf("Sender must be specified")
	}
	return message, nil
}
