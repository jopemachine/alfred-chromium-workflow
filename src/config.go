package src

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/deanishe/awgo"
)

type WorkflowConfig struct {
	Browser            string
	Locale             string
	Profile            string
	SwitchableProfiles string
	ResultCountLimit   int
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
		Arg(browserName).
		Icon(&aw.Icon{fmt.Sprintf(`assets/browser-icons/%s.png`, browserName), ""}).
		Autocomplete(browserName)
}

var SelectBrowser = func(wf *aw.Workflow, query string) {
	browsers := []string{
		"Chrome",
		"Chrome Canary",
		"Chromium",
		"Edge",
		"Brave",
		"Naver Whale",
		"Epic",
	}

	for _, browser := range browsers {
		if CheckBrowserIsInstalled(browser) {
			addNewBrowserItem(wf, browser)
		}
	}

	wf.Filter(query)
}

var ChangeBrowser = func(browserName string) {
	err := ConfigAPI.Set("BROWSER", browserName, true).Do()
	CheckError(err)

	// Change workflow's default icon
	CopyFile(fmt.Sprintf("assets/browser-icons/%s.png", browserName), "icon.png")

	fmt.Print(browserName)
}

var SelectProfile = func(wf *aw.Workflow, query string) {
	profileRoot := GetProfileRootPath(Conf.Browser)
	profileFilePaths, err := filepath.Glob(profileRoot + "/" + "Profile *")
	CheckError(err)

	var profiles []string

	for _, profileFilePath := range profileFilePaths {
		profileFilePathArr := strings.Split(profileFilePath, "/")
		profiles = append(profiles, profileFilePathArr[len(profileFilePathArr)-1])
	}

	possibleProfiles := strings.Split(Conf.SwitchableProfiles, ",")
	possibleProfiles = append(possibleProfiles, profiles...)

	for _, profile := range possibleProfiles {
		wf.NewItem(profile).
			Valid(true).
			Arg(profile).
			Autocomplete(profile)
	}

	wf.Filter(query)
}

var ChangeProfile = func(profileName string) {
	// Check if the profile folder exist in the browser first
	if !FileExist(GetProfileRootPath(Conf.Browser)) {
		fmt.Print("")
		return
	}

	err := ConfigAPI.Set("PROFILE", profileName, true).Do()
	CheckError(err)
	fmt.Print(profileName)
}

var GetFaviconDirectoryPath = func(wf *aw.Workflow) string {
	// Make cache object to ensure the path exists
	faviconCache := aw.NewCache(filepath.Join(wf.Cache.Dir, "favicon"))
	return faviconCache.Dir
}

var GetTempDataPath = func(wf *aw.Workflow) string {
	temp := aw.NewCache(filepath.Join(wf.Data.Dir, "temp"))
	return temp.Dir
}
