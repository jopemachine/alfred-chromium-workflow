package src

// Ref: https://github.com/bit2pixel/chrome-control

import (
	"fmt"
	"strings"
	"encoding/json"

	"github.com/deanishe/awgo/util"
)

var getApplicationName = func (browserName string) string {
	switch browserName {
	case "Edge":
		return "Microsoft Edge"
	case "Chromium":
		return "Google Chrome"
	case "Chrome":
		fallthrough
	case "Chrome Canary":
		return "Google Chrome"
	default:
		panic("Unsupported browser. Please make a issue to support the browser if the browser is based on Chromium.")
	}
}

var getListUpTabScript = func () string {
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

					if (!query || title.includes(query) || url.includes(query)) {
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

var getCloseTabScript = func () string {
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

var getFocusTabScript = func () string {
	browserName := getApplicationName(Conf.Browser)

	return fmt.Sprintf(`
	const browser = Application('%s')
	browser.includeStandardAdditions = true

	function focusTab(winIdx, tabIdx) {
		browser.windows[winIdx].visible = true
		browser.windows[winIdx].activeTabIndex = tabIdx + 1 // Focus on tab
		browser.windows[winIdx].index = 1 // Focus on this specific browser window
		browser.activate()
	}

	function run(argv) {
		return focusTab(Number(argv[0]), Number(argv[1]))
	}
	`, browserName)
}

var ListOpenedTabs = func (query string) {
	stdout, err := util.RunJS(getListUpTabScript(), query)
	CheckError(err)

	var serializedStdout map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &serializedStdout)
	CheckError(err)

	for _, item := range serializedStdout["items"].([]interface{}) {
		url := item.(map[string]interface{})["url"].(string)
		domainName := ExtractDomainName(url)
		iconPath := fmt.Sprintf(`cache/%s.png`, domainName)

		if FileExist(iconPath) {
			item.(map[string]interface{})["icon"] = map[string]string { "path": iconPath }
		}
	}

	result, err := json.Marshal(serializedStdout)
	CheckError(err)
	fmt.Print(string(result))
}

var CloseTab = func (query string) {
	argv := strings.Split(query, ",")
	_, err := util.RunJS(getCloseTabScript(), argv...)
	CheckError(err)
}

var FocusTab = func (query string) {
	argv := strings.Split(query, ",")
	_, err := util.RunJS(getFocusTabScript(), argv...)
	CheckError(err)
}
