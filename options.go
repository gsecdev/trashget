package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Port       int    `short:"p" long:"port" default:"8000" description:"port to listen at"`
	IP         string `short:"i" long:"ip" description:"IP to listen at (defaults to all IPs)"`
	Filename   string `short:"f" long:"filename" default:"full_backup.zip" description:"filename to serve"`
	Size       int64  `short:"s" long:"size" default:"1000" description:"virtual size of file (in MB)"`
	Uri        string `short:"u" long:"uri" default:"/" description:"URI to serve at"`
	Throttle   int    `short:"t" long:"throttle" default:"-1" description:"throttle bandwith (in Mbit/s)"`
	AbortAfter int    `short:"a" long:"abortAfter" default:"-1" description:"abort transmission after given %"`
}

func (o Options) DoesThrottle() bool {
	return o.Throttle > 0
}

func (o Options) DoesAbort() bool {
	return cmdOpts.AbortAfter != -1
}

func (o *Options) Validate(writeHelp func()) {
	if cmdOpts.Filename == "" {
		writeHelp()
		log.Fatalf("empty filename specified")
	}

	if cmdOpts.Port < 0 || cmdOpts.Port > 65535 {
		writeHelp()
		log.Fatalf("port needs to be in rage 0-65535")
	}

	if cmdOpts.Size < 0 {
		writeHelp()
		log.Fatalf("illegal file size: %d", cmdOpts.Size)
	}

	if (cmdOpts.AbortAfter < 0 || cmdOpts.AbortAfter > 100) && cmdOpts.AbortAfter != -1 {
		writeHelp()
		log.Fatalf("illegal abort percentage: %d. needs to be -1 (deactivated) or between 0-100", cmdOpts.AbortAfter)
	}
}

func (o *Options) parseFlags() (err error) {
	parser := flags.NewParser(o, flags.HelpFlag|flags.PassDoubleDash)

	writeHelp := func() {
		parser.WriteHelp(os.Stdout)
	}

	if _, err = parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			writeHelp()
			os.Exit(0)
		} else {
			writeHelp()
			return fmt.Errorf("error parsing flags: %v", err)
		}
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.00",
	})

	o.Validate(writeHelp)

	// log.SetOutput(os.Stdout)

	// verbosity := cmdOpts.GetVerbosity()

	// if verbosity.Debug {
	// 	log.SetLevel(log.DebugLevel)
	// } else if verbosity.Quiet {
	// 	log.SetLevel(log.WarnLevel)
	// } else {
	// 	log.SetLevel(log.InfoLevel)
	// }
	return nil
}
