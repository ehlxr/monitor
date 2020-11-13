package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ehlxr/monitor/pkg"
	"github.com/hpcloud/tail"
	log "unknwon.dev/clog/v2"

	"regexp"

	dt "github.com/JetBlink/dingtalk-notify-go-sdk"
)

var (
	dingTalk *dt.Robot
	//limiter  *pkg.LimiterServer
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
	//limiter = pkg.NewLimiterServer(1*time.Minute, 20)

	tailFile()
}

func sendMsg(content string) {
	if err := dingTalk.SendMarkdownMessage(
		"new message",
		fmt.Sprintf("%s\n%s", pkg.Opts.AppName, content),
		pkg.Opts.Robot.AtMobiles,
		pkg.Opts.Robot.IsAtAll,
	); err != nil {
		log.Error("%+v", err)

		return
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

	var buffer strings.Builder
	var times int
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			<-ticker.C
			//if buffer.Len() > 0 && times > 2 {
			if buffer.Len() > 0 {
				sendMsg(buffer.String())
				buffer.Reset()
			}

			//buffer.Reset()
			//times = 0
		}
	}()

	for line := range tf.Lines {
		text := line.Text
		if pkg.Opts.KeyWordIgnoreCase {
			text = strings.ToLower(text)
		}

		keys := strings.Split(pkg.Opts.KeyWord, ",")
		for _, key := range keys {
			if ok, _ := regexp.Match(strings.TrimSpace(key), []byte(text)); ok {
				//if limiter.IsAvailable() {
				//	sendMsg("- " + line.Text + "\n")
				//} else {
				//	log.Error("dingTalk 1 m allow send 20 msg. msg %v discarded.",
				//		line.Text)
				//}
				buffer.WriteString("- " + line.Text + "\n")
				times++

				break
			}
		}
	}
}
