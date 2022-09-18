package src

// Implement tab, window related features using JXA, Applescript.
// Ref: https://github.com/bit2pixel/chrome-control

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
)

var getApplicationName = func(browserName string) string {
	switch browserName {
	case "Opera":
		return "Opera"
	case "Edge":
		return "Microsoft Edge"
	case "Edge Canary":
		return "Microsoft Edge Canary"
	case "Chrome":
		return "Google Chrome"
	case "Chrome Canary":
		return "Google Chrome Canary"
	case "Chrome Beta":
		return "Google Chrome Beta"
	case "Brave":
		return "Brave Browser"
	case "Chromium":
		return "Chromium"
	case "Naver Whale":
		return "Whale"
	case "Epic":
		return "Epic"
	case "Vivaldi":
		return "Vivaldi"
	default:
		panic("Unsupported browser. Please consider to make a issue to support the browser if the browser is based on Chromium.")
	}
}

var getListUpTabScript = func() string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	const browser = Application('%s')
	browser.includeStandardAdditions = true

	function list(query) {
		 let allTabsTitle = browser.windows.tabs.title()
		 let allTabsUrls = browser.windows.tabs.url()

		 var titleToUrl = {}
		 for (var winIdx = 0; winIdx < allTabsTitle.length; winIdx++) {
			  for (var tabIdx = 0; tabIdx < allTabsTitle[winIdx].length; tabIdx++) {
					let title = allTabsTitle[winIdx][tabIdx]
					let url = allTabsUrls[winIdx][tabIdx]

					if (!query || title.toLowerCase().includes(query) || url.toLowerCase().includes(query)) {
						titleToUrl[title] = {
							 'title': title || 'No Title',
							 'url': url,
							 'winIdx': winIdx,
							 'tabIdx': tabIdx,
							 'arg': winIdx + ',' + tabIdx,
							 'subtitle': url,
						}
					}
			  }
		 }

		 out = { 'items': [] }
		 Object.keys(titleToUrl).sort().forEach(title => {
			  out.items.push(titleToUrl[title])
		 })

		 return JSON.stringify(out)
	}

	function run(argv) {
		return list(argv)
	}
	`, browserName)
}

var getCloseTabScript = func() string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	const browser = Application('%s')
	browser.includeStandardAdditions = true

	function closeTab(winIdx, tabIdx) {
		let tabToClose = browser.windows[winIdx].tabs[tabIdx]
		tabToClose.close()
	}

	function run(argv) {
		return closeTab(Number(argv[0]), Number(argv[1]))
	}
	`, browserName)
}

var getFocusTabScript = func() string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	const browser = Application('%s')
	browser.includeStandardAdditions = true

	function focusTab(winIdx, tabIdx) {
		browser.windows[winIdx].visible = true
		browser.windows[winIdx].activeTabIndex = tabIdx + 1
		browser.windows[winIdx].index = 1
		browser.activate()
	}

	function run(argv) {
		return focusTab(Number(argv[0]), Number(argv[1]))
	}
	`, browserName)
}

var getNewWindowScript = func() string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	do shell script "open -na '%s' --args --profile-directory='%s'"
	`, browserName, Conf.Profile)
}

var getNewTabScript = func() string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	tell application "%s"
		activate
		tell front window to make new tab at after (get active tab)
	end tell
	`, browserName)
}

var ListOpenedTabs = func(wf *aw.Workflow, query string) {
	query = strings.ToLower(query)
	stdout, err := util.RunJS(getListUpTabScript(), query)
	CheckError(err)

	var serializedStdout map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &serializedStdout)
	CheckError(err)
	items := serializedStdout["items"].([]interface{})

	for _, item := range items {
		url := item.(map[string]interface{})["url"].(string)
		domainName := ExtractDomainName(url)
		item.(map[string]interface{})["domain"] = domainName

		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)
		if FileExist(iconPath) {
			item.(map[string]interface{})["icon"] = map[string]string{"path": iconPath}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		ithDomainName := items[i].(map[string]interface{})["domain"].(string)
		jthDomainName := items[j].(map[string]interface{})["domain"].(string)

		return ithDomainName > jthDomainName
	})

	result, err := json.Marshal(serializedStdout)
	CheckError(err)
	fmt.Print(string(result))
}

var CloseTab = func(query string) {
	argv := strings.Split(query, ",")
	_, err := util.RunJS(getCloseTabScript(), argv...)
	CheckError(err)
}

var FocusTab = func(query string) {
	argv := strings.Split(query, ",")
	_, err := util.RunJS(getFocusTabScript(), argv...)
	CheckError(err)
}

var OpenNewTab = func() {
	_, err := util.RunAS(getNewTabScript())
	// Open New Window instead
	if err != nil {
		OpenNewWindow()
	}
}

var OpenNewWindow = func() {
	_, err := util.RunAS(getNewWindowScript())
	CheckError(err)
}
