package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Version bool   `short:"v" long:"version" description:"Show version"`
	Url     string `long:"url" description:"Database connection string"`
	Host    string `short:"h" long:"host" description:"Server hostname or IP Address" default:"localhost"`
	Port    uint16 `short:"p" long:"port" description:"Server port" default:"5432"`
	User    string `short:"u" long:"user" description:"Database user"`
	Pass    string `long:"pass" description:"User Password"`
	DBName  string `long:"db" description:"Database name" default:"postgres"`
	Ssl     string `long:"ssl" description:"SSL option" default:"disable"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("pgmonit v%s\n", Version)
		os.Exit(0)
	}

	if opts.Url == "" {
		if opts.Url = os.Getenv("DATABASE_URL"); opts.Url == "" {
			if opts.User == "" {
				fmt.Fprintf(os.Stderr, "user required")
				os.Exit(0)
			}

			opts.Url = fmt.Sprintf(
				"host=%s port=%d user=%s dbname=%s sslmode=%s",
				opts.Host, opts.Port,
				opts.User, opts.DBName, opts.Ssl,
			)

			if opts.Pass != "" {
				opts.Url += fmt.Sprintf(" password=%s", opts.Pass)
			}
		}
	}

	DB, err = NewDB(opts.Url)
	if err != nil {
		fmt.Printf("DB connection error: %s\n", err.Error())
		os.Exit(1)
	}
	defer DB.Close()

	AppUp()
}
