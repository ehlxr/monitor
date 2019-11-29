package main

import (
	"fmt"
	"github.com/ehlxr/monitor/pkg"
	"github.com/hpcloud/tail"
	"strings"
	"time"
	log "unknwon.dev/clog/v2"

	dt "github.com/JetBlink/dingtalk-notify-go-sdk"
	"regexp"
)

var (
	dingTalk *dt.Robot
	limiter  *pkg.LimiterServer
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func main() {
	pkg.ParseArg()

	dingTalk = dt.NewRobot(pkg.Opts.Robot.Token, pkg.Opts.Robot.Secret)
	limiter = pkg.NewLimiterServer(1*time.Minute, 20)

	tailFile()
}

func sendMsg(content string) {
	if err := dingTalk.SendTextMessage(
		fmt.Sprintf("%s\n%s", pkg.Opts.AppName, content),
		pkg.Opts.Robot.AtMobiles,
		pkg.Opts.Robot.IsAtAll,
	); err != nil {
		log.Error("%+v", err)
	}

	log.Info("send message <%s> success", content)
}

func tailFile() {
	tf, err := tail.TailFile(pkg.Opts.File,
		tail.Config{
			ReOpen:   true,
			Follow:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		})
	if err != nil {
		log.Fatal("Tail file %+v", err)
	}

	if pkg.Opts.KeyWordIgnoreCase {
		pkg.Opts.KeyWord = strings.ToLower(pkg.Opts.KeyWord)
	}

	log.Info("monitor app <%s> file <%s>, filter by <%s>, ignore case <%v>...",
		pkg.Opts.AppName,
		pkg.Opts.File,
		pkg.Opts.KeyWord,
		pkg.Opts.KeyWordIgnoreCase)

	for line := range tf.Lines {
		text := line.Text
		if pkg.Opts.KeyWordIgnoreCase {
			text = strings.ToLower(text)
		}

		keys := strings.Split(pkg.Opts.KeyWord, ",")
		for _, key := range keys {
			if ok, _ := regexp.Match(strings.TrimSpace(key), []byte(text)); ok {
				if limiter.IsAvailable() {
					sendMsg(line.Text)
				} else {
					log.Error("dingTalk 1 m allow send 20 msg. msg %v discarded.",
						line.Text)
				}
				break
			}
		}
	}
}
