package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/doublecloud/kfbench/internal/franz"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "kfbench",
		Commands: []*cli.Command{
			franz.Command(),
		},
	}


	// Expose :8080/debug/pprof
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	
}
