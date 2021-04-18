package main

/*
	gouploadserver [options] [path]
	ex: gouploadserver -p 8081 .
	options:
		-p or --port Port to use (defaults to 8080)

*/

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/guilhermerodrigues680/gouploadserver/app"

	"github.com/sirupsen/logrus"
)

var portFlag = flag.Int("port", 8000, "Port to use")
var watchMemUsageFlag = flag.Bool("watch-mem", false, "Watch memory usage")
var devFlag = flag.Bool("dev", false, "Use development settings")
var keepOriginalUploadFileNameFlag = flag.Bool("keep-upload-filename", false, "Keep original upload file name: Use 'filename.ext' instead of 'filename<-random>.ext'")
var pathArg string

func main() {
	// usage: flag -h or --help
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gouploadserver [options] [path]")
		fmt.Fprintln(flag.CommandLine.Output(), "[path] defaults to ./")
		fmt.Fprintln(flag.CommandLine.Output(), "Options are:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(flag.CommandLine.Output(), "  --%-24v %v (default %v)\n", f.Name, f.Usage, f.DefValue)
		})
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Powered By: guilhermerodrigues680")
	}

	// parses the command-line flags
	flag.Parse()

	logger := getLogger(*devFlag)
	logger.Trace(strings.Join(os.Args, " "))

	if *devFlag {
		flag.VisitAll(func(f *flag.Flag) {
			logger.Debugf("--%v (value %v) (default %v)", f.Name, f.Value, f.DefValue)
		})
	}

	if *watchMemUsageFlag {
		go func() {
			for {
				time.Sleep(time.Second)
				PrintMemUsage(logger.WithField("log", "memstats"))
			}
		}()
	}

	wd := flag.Arg(0)
	if wd == "" {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
		}
		wd = cwd
	}

	err := app.Run(wd, *portFlag, *keepOriginalUploadFileNameFlag, logger.WithField("app", "run"))
	if err != nil {
		logger.Fatal(err)
	}
}

func getLogger(development bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: development,
	})

	if development {
		logger.SetLevel(logrus.TraceLevel) // log all
	} else {
		logger.SetLevel(logrus.InfoLevel) // log only info and above
	}

	logger.SetOutput(os.Stdout) // Output to stdout instead of the default stderr
	return logger
}

func PrintMemUsage(logger *logrus.Entry) {
	// For info, see: https://golang.org/pkg/runtime/#MemStats
	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.Infof("Alloc = %v MiB\tHeapAlloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v",
		bToMb(m.Alloc), bToMb(m.HeapAlloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}