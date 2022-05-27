package src

import (
	"fmt"
	"strings"

	"github.com/deanishe/awgo"
)

type WorkflowConfig struct {
	Browser            string
	Locale             string
	Profile            string
	SwitchableProfiles string
	ResultLimitCount   uint8
}

var wf *aw.Workflow

var Conf = &WorkflowConfig{}

var ConfigAPI *aw.Config

var ImportConfig = func() {
	ConfigAPI = aw.NewConfig()

	if err := ConfigAPI.To(Conf); err != nil {
		panic(err)
	}
}

func addNewBrowserItem(wf *aw.Workflow, browserName string) {
	wf.NewItem(browserName).
		Valid(true).
		Arg(browserName)
}

var SelectBrowser = func(wf *aw.Workflow, query string) {
	addNewBrowserItem(wf, "Chrome")
	addNewBrowserItem(wf, "Chrome Canary")
	addNewBrowserItem(wf, "Chromium")
	addNewBrowserItem(wf, "Edge")
	addNewBrowserItem(wf, "Brave")

	wf.Filter(query)
}

var ChangeBrowser = func(browserName string) {
	err := ConfigAPI.Set("BROWSER", browserName, true).Do()
	CheckError(err)
	fmt.Print(browserName)
}

var SelectProfile = func(wf *aw.Workflow, query string) {
	for _, profile := range strings.Split(Conf.SwitchableProfiles, ",") {
		wf.NewItem(profile).
			Valid(true).
			Arg(profile)
	}

	wf.Filter(query)
}

var ChangeProfile = func(profileName string) {
	err := ConfigAPI.Set("PROFILE", profileName, true).Do()
	CheckError(err)
	fmt.Print(profileName)
}
