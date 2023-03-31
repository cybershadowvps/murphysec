package env

import (
	"os"
	"strconv"
	"strings"
)

func envi(name string, defaultValue int) int {
	if i, e := strconv.Atoi(os.Getenv(name)); e != nil {
		return defaultValue
	} else {
		return i
	}
}

var SpecificProjectName = ""
var DisableGit = false
var _ServerBaseURL = "https://www.murphysec.com"

var DisableMvnCommand = strings.TrimSpace(os.Getenv("NO_MVN")) != ""
var MavenCentral string
var TlsAllowInsecure = os.Getenv("TLS_ALLOW_INSECURE") != ""

var SkipGradleExecution = os.Getenv("SKIP_GRADLE_EXECUTION") != ""

func init() {
	if strings.TrimSpace(os.Getenv("SKIP_MAVEN_CENTRAL")) == "" {
		MavenCentral = "https://repo1.maven.org/maven2/"
	}
}

func ConfigureServerBaseUrl(u string) {
	_ServerBaseURL = strings.TrimRight(strings.TrimSpace(u), "/")
}

func ServerBaseUrl() string {
	return _ServerBaseURL
}
