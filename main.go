package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/SouleBA/httpmon/monitor"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var treshold uint
	var filePath string

	flag.UintVar(&treshold, "treshold", 10, "treshold /s after which to be alerted.")
	flag.StringVar(&filePath, "filePath", "/tmp/access.log", "path to the file to monitor.")

	flag.Parse()

	l := monitor.NewLauncher(monitor.DefaultConfig(filePath, treshold))
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	writer.Flush()
	go l.Launch(writer)

	<-sigs
	l.Shutdown()

}
