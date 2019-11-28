package pkg

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var (
	Opts struct {
		AppName           string `short:"n" long:"monitor-app-name" env:"MONITOR_APP_NAME" description:"The name of the application being monitored, which will be added to the content before"`
		File              string `short:"f" long:"monitor-file" env:"MONITOR_FILE" description:"The file to be monitored" required:"true"`
		KeyWord           string `short:"k" long:"search-keyword" env:"SEARCH_KEYWORD" description:"Keyword to be search for, Multiple values separated by ','" default:"ERRO"`
		KeyWordIgnoreCase bool   `short:"c" long:"keyword-case-sensitive" env:"KEYWORD_IGNORE_CASE" description:"Whether Keyword ignore case"`
		Version           bool   `short:"v" long:"version" description:"Show version info"`
		Robot             Robot  `group:"DingTalk Robot Options" namespace:"robot" env-namespace:"ROBOT" `
	}
)

type Robot struct {
	Token     string   `short:"t" long:"token" env:"TOKEN" description:"DingTalk robot access token" required:"true"`
	Secret    string   `short:"s" long:"secret" env:"SECRET" description:"DingTalk robot secret"`
	AtMobiles []string `short:"m" long:"at-mobiles" env:"AT_MOBILES" env-delim:"," description:"The mobile of the person will be at"`
	IsAtAll   bool     `short:"a" long:"at-all" env:"AT_ALL" description:"Whether at everyone"`
}

func ParseArg() {
	parser := flags.NewParser(&Opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	if AppName != "" {
		parser.Name = AppName
	}

	if _, err := parser.Parse(); err != nil {
		if Opts.Version {
			// -v
			PrintVersion()
			os.Exit(0)
		}

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			// -h
			_, _ = fmt.Fprintln(os.Stdout, err)
			os.Exit(0)
		}

		// err
		_, _ = fmt.Fprintln(os.Stderr, err)

		parser.WriteHelp(os.Stderr)

		os.Exit(1)
	}
}
