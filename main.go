package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	log "unknwon.dev/clog/v2"

	"net/http"
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
	bannerBase64 = "DQogX19fXyAgX19fXyAgICBfX18gIF9fX19fIA0KKCAgXyBcKCAgXyBcICAvIF9fKSggIF8gICkNCiApKF8pICkpKF8pICkoIChfLS4gKShfKSggDQooX19fXy8oX19fXy8gIFxfX18vKF9fX19fKQ0K"

	opts struct {
		MonitorFile string `short:"m" long:"monitor-file" env:"MONITOR_FILE" description:"The file to be monitored" required:"true"`
		KeyWord     string `short:"k" long:"key-word" env:"KEY_WORD" description:"Key word to be filter" required:"true"`
		WebHookUrl  string `short:"u" long:"webhook-url" env:"URL" description:"Webhook url of dingtalk" required:"true"`
		Version     bool   `short:"v" long:"version" description:"Show version info"`
	}
)

func init() {
	initLog()
}

func main() {
	parseArg()

	tf, err := tail.TailFile(opts.MonitorFile,
		tail.Config{
			Follow:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		})
	if err != nil {
		log.Fatal("Tail file %+v", err)
	}

	for line := range tf.Lines {
		if ok, _ := regexp.Match(opts.KeyWord, []byte(line.Text)); ok {
			log.Info("%s", dingToInfo(line.Text))
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
	if AppName != "" {
		parser.Name = AppName
	}

	if _, err := parser.Parse(); err != nil {
		if opts.Version {
			printVersion()
			os.Exit(0)
		}

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			_, _ = fmt.Fprintln(os.Stdout, err)
			os.Exit(0)
		}

		_, _ = fmt.Fprintln(os.Stderr, err)

		parser.WriteHelp(os.Stderr)

		os.Exit(1)
	}
}

func dingToInfo(msg string) []byte {
	content, data := make(map[string]string), make(map[string]interface{})

	content["content"] = msg
	data["msgtype"] = "text"
	data["text"] = content
	b, _ := json.Marshal(data)

	log.Info("send to %s data <%s>",
		opts.WebHookUrl,
		b)

	resp, err := http.Post(opts.WebHookUrl,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		log.Error("send request to %s %+v",
			opts.WebHookUrl,
			err)

	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("send to %s data <%s> result is %s",
		opts.WebHookUrl,
		b,
		body)
	return body
}

// printVersion Print out version information
func printVersion() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, AppName, Version, BuildTime, GitCommit, GoVersion)
}
