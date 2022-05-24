package src

import (
	"os"

	"github.com/deanishe/awgo"
)

type WorkflowConfig struct {
	Browser string
	Locale string
	ChromeProfile string
	ExcludeDomains []string
	ResultLimitCount uint8
}

var wf *aw.Workflow

var Conf = &WorkflowConfig{}

var InitConfig = func () {
	os.Setenv("BROWSER", "Chrome")
	os.Setenv("LOCALE", "en_US")
	os.Setenv("CHROME_PROFILE", "Profile 3")
	os.Setenv("EXCLUDE_DOMAINS", "")
	os.Setenv("RESULT_LIMIT_COUNT", "5")
}

var ImportConfig = func () {
	InitConfig()
	conf := aw.NewConfig()

	if err := conf.To(Conf); err != nil {
		panic(err)
	}
}

