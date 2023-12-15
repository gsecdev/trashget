package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Port     int    `short:"p" long:"port" description:"port to listen at"`
	IP       string `short:"i" long:"ip" description:"IP to listen at (defaults to all IPs)"`
	Filename string `short:"f" long:"filename" default:"full_backup.zip" description:"filename to serve"`
	Size     int64  `short:"s" long:"size" description:"virtual size to server (in MB)"`
	Uri      string `short:"u" long:"uri" default:"/" description:"URI to serve at"`
	// Throttle float32 `short:"t" long:"throttle" default:"-1" description:"throttle bandwith (in Mbit/s)"`
}

// func (o Options) DoThrottle() bool {
// 	return o.Throttle > 0
// }

func (o *Options) parseFlags() (writeHelp func(), err error) {
	parser := flags.NewParser(o, flags.HelpFlag|flags.PassDoubleDash)

	writeHelp = func() {
		parser.WriteHelp(os.Stdout)
	}

	if _, err = parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			os.Exit(0)
		} else {
			parser.WriteHelp(os.Stdout)
			return writeHelp, fmt.Errorf("error parsing flags: %v", err)
		}
	}

	// log.SetFormatter(&log.TextFormatter{
	// 	ForceColors:     false,
	// 	FullTimestamp:   true,
	// 	TimestampFormat: "060102 150405.00",
	// })

	// log.SetOutput(os.Stdout)

	// verbosity := cmdOpts.GetVerbosity()

	// if verbosity.Debug {
	// 	log.SetLevel(log.DebugLevel)
	// } else if verbosity.Quiet {
	// 	log.SetLevel(log.WarnLevel)
	// } else {
	// 	log.SetLevel(log.InfoLevel)
	// }
	return writeHelp, nil
}
