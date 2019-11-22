package pkg

import (
	"encoding/base64"
	"fmt"
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
)

// PrintVersion Print out version information
func PrintVersion() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, AppName, Version, BuildTime, GitCommit, GoVersion)
}
