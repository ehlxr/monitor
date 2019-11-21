package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
	"strings"
	log "unknwon.dev/clog/v2"

	dt "github.com/JetBlink/dingtalk-notify-go-sdk"
	"os"
	"regexp"
)

var (
	AppName   string
	Version   string
	BuildTime string
	GitCommit string
	GoVersion string

	versionTpl = `%s
Name: %s
Version: %s
BuildTime: %s
GitCommit: %s
GoVersion: %s

`
	bannerBase64 = "DQogX18gIF9fICBfX19fXyAgXyAgXyAgX19fXyAgX19fXyAgX19fX18gIF9fX18gDQooICBcLyAgKSggIF8gICkoIFwoICkoXyAgXykoXyAgXykoICBfICApKCAgXyBcDQogKSAgICAoICApKF8pKCAgKSAgKCAgXykoXyAgICkoICAgKShfKSggICkgICAvDQooXy9cL1xfKShfX19fXykoXylcXykoX19fXykgKF9fKSAoX19fX18pKF8pXF8pDQo="

	opts struct {
		File    string `short:"f" long:"monitor-file" env:"MONITOR_FILE" description:"The file to be monitored" required:"true"`
		KeyWord string `short:"k" long:"search-keyword" env:"SEARCH_KEYWORD" description:"Key word to be search for" default:"ERRO"`
		Version bool   `short:"v" long:"version" description:"Show version info"`
		Robot   robot  `group:"DingTalk Robot Options" namespace:"robot" env-namespace:"ROBOT" `
	}
)

type robot struct {
	Token     string   `short:"t" long:"token" env:"TOKEN" description:"DingTalk robot access token" required:"true"`
	Secret    string   `short:"s" long:"secret" env:"SECRET" description:"DingTalk robot secret"`
	AtMobiles []string `short:"m" long:"at-mobiles" env:"AT_MOBILES" env-delim:"," description:"The mobile of the person will be at"`
	IsAtAll   bool     `short:"a" long:"at-all" env:"AT_ALL" description:"Whether at everyone"`
}

func init() {
	initLog()
}

func main() {
	parseArg()

	tf, err := tail.TailFile(opts.File,
		tail.Config{
			ReOpen:   true,
			Follow:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		})
	if err != nil {
		log.Fatal("Tail file %+v", err)
	}
	log.Info("monitor file %s...", opts.File)

	dingTalk := dt.NewRobot(opts.Robot.Token, opts.Robot.Secret)

	opts.KeyWord = strings.ToLower(opts.KeyWord)
	for line := range tf.Lines {
		text := strings.ToLower(line.Text)
		if ok, _ := regexp.Match(opts.KeyWord, []byte(text)); ok {
			err = dingTalk.SendTextMessage(line.Text, opts.Robot.AtMobiles, opts.Robot.IsAtAll)
			if err != nil {
				log.Error("%+v", err)
				continue
			}

			log.Info("send message <%s> success", line.Text)
		}
	}
}

func initLog() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func parseArg() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	if AppName != "" {
		parser.Name = AppName
	}

	if _, err := parser.Parse(); err != nil {
		if opts.Version {
			// -v
			printVersion()
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

// printVersion Print out version information
func printVersion() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, AppName, Version, BuildTime, GitCommit, GoVersion)
}
