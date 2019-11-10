// Copyright (c) 2019, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.

package sendmail

import (
	"fmt"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	BodyType          string   `short:"B" description:"(ignored) Set the body type to type. Current legal values are 7BIT or 8BITMIME."`
	BParam            []string `short:"b"`
	DefaultDelivery   bool     `short:"bm" description:"Deliver mail in the usual way (default)." no-flag:"true"`
	PrintQueue        bool     `short:"bp" description:"(ignored) Print a listing of the queue(s)." no-flag:"true"`
	UseSMTP           bool     `short:"bs" description:"(ignored) Use the SMTP protocol as described in RFC821 on standard input and output." no-flag:"true"`
	ConfigurationFile string   `short:"C" description:"(ignored) Use alternate configuration file."`
	DebuggingFlag     string   `short:"d" description:"(ignored) Set the debugging flag for category to level."`
	SenderFullName    string   `short:"F" description:"Set the sender full name."`
	SenderAddress     string   `short:"f" description:"Set the envelope sender address."`
	HopCount          string   `short:"h" description:"(ignored) Set the hop count to N."`
	IgnoreDot         bool     `short:"i" description:"When reading a message from standard input, don't treat a line with only a . character as the end of input."`
	SyslogTag         string   `short:"L" description:"(ignored) Set the identifier used in syslog messages to the supplied tag."`
	FlagM             string   `short:"m" description:"(ignored)"`
	DeliveryStatusNot string   `short:"N" description:"(ignored) Set delivery status notification conditions."`
	NotAliasing       string   `short:"n" description:"(ignored) Don't do aliasing."`
	OptionsLong       []string `short:"O" description:"(ignored) Set sendmail option [option=value] to the specified value."`
	Options           []string `short:"o" description:"(ignored) Set sendmail option [x] to the specified value."`
	EParam            []string `short:"e"`
	FlagEM            bool     `short:"em" description:"(ignored)" no-flag:"true"`
	FlagEP            bool     `short:"ep" description:"(ignored)" no-flag:"true"`
	FlagEQ            bool     `short:"eq" description:"(ignored)" no-flag:"true"`
	FlagP             string   `short:"p" description:"(ignored)"`
	FlagQ             string   `short:"q" description:"(ignored)"`
	FlagR             string   `short:"R" description:"(ignored)"`
	LegacyFrom        string   `short:"r" description:"An alternate and obsolete form of the -f flag."`
	ExtractRecipients bool     `short:"t" description:"Extract recipients from message headers. These are added to any recipients specified on the command line."`
	FlagU             string   `short:"U" description:"(ignored)"`
	EnvelopeID        string   `short:"V" description:"(ignored) Set the original envelope id."`
	Verbose           bool     `short:"v" description:"(ignored) Go into verbose mode."`
	LogFile           string   `short:"X" description:"(ignored) Log all traffic in and out of mailers in the indicated log file.\n"`
	Version           bool     `long:"version" hidden:"true"`
	ShowHelp          bool     `long:"help" description:"Display this help and exit"`
}

func showHelp(data interface{}) {
	t := reflect.TypeOf(data).Elem()

	fmt.Print("usage: sendmail [options ...] [recipient ...]\n" +
		"A sendmail compatible interface for sentry events.\n\n" +
		"options:\n")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		usage := field.Tag.Get("description")
		short := field.Tag.Get("short")
		long := field.Tag.Get("long")

		if usage != "" {
			if short != "" {
				fmt.Printf("  -%s\t  %s\n", short, usage)
			}
			if long != "" {
				fmt.Printf("  --%s  %s\n", long, usage)
			}
		}
	}

	fmt.Print("\n")
}

func ParseOptions() {
	// To handle flags as sendmail we need to support
	// GNU extensions to the POSIX recommendations for command-line options
	parser := flags.NewParser(&opts, flags.None)
	args, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if opts.ShowHelp {
		// Display custom help and exit
		showHelp(&opts)
		os.Exit(0)
	}

	if opts.Version {
		fmt.Printf("v%s\n", Version)
		os.Exit(0)
	}

	if len(args) > 0 {
		recipients = args[0]
	}

	// TODO: Process BParams if we want to support -bd
}
