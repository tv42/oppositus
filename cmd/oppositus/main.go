package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"eagain.net/go/oppositus"
	"eagain.net/go/oppositus/internal/config"
	"eagain.net/go/oppositus/internal/version"
	"golang.org/x/net/context"
)

var (
	showVersion = flag.Bool("version", false, "display version and exit")
)

func doit(configPath string, dest string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// jump through hoops to clean up temp files on control-C
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	defer close(signals)
	defer signal.Stop(signals)
	go func() {
		<-signals
		cancel()
	}()

	conf, err := config.Load(configPath)
	if err != nil {
		return err
	}
	success := true
	errFn := func(err error) error {
		log.Printf("%v", err)
		success = false
		return nil
	}
	opts := []oppositus.Option{
		oppositus.WithFilter(conf.Filters.Match),
		oppositus.WithErrorHandler(errFn),
	}
	if conf.Channels != nil {
		opts = append(opts, oppositus.WithChannels(conf.Channels...))
	}
	if err := oppositus.Mirror(ctx, dest, opts...); err != nil {
		return err
	}
	if !success {
		return errors.New("mirror failed")
	}
	return nil
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [OPTS] CONFIG DEST\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s %s\n", prog, version.Version)
		os.Exit(0)
	}
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	configPath := flag.Arg(0)
	dest := flag.Arg(1)

	if err := doit(configPath, dest); err != nil {
		log.Fatal(err)
	}
}
