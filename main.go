package main

import (
	"flag"

	"github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	API "github.com/jopemachine/alfred-chromium-workflow/src"
	"github.com/klauspost/lctime"
)

var wf *aw.Workflow

func alfredCallback() {
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		commandType := args[0]

		var query string

		if len(args) > 1 {
			query = args[1]
		} else {
			query = ""
		}

		API.ShowUpdateStatus(wf, query)

		switch commandType {
		// Fetch data
		case "search-log":
			API.FetchSearchData(wf, query)
		case "visit-history":
			API.FetchHistory(wf, query)
		case "bookmark":
			API.FetchBookmark(wf, query)
		case "bookmark-folder":
			API.FetchBookmarkFolder(wf, query)
		case "download":
			API.FetchDownloadHistory(wf, query)
		case "login":
			API.FetchLoginData(wf, query)
		case "autofill":
			API.FetchAutofillData(wf, query)

		// Tab, window related features
		case "listup-tabs":
			API.ListOpenedTabs(query)
			return
		case "close-tab":
			API.CloseTab(query)
			return
		case "focus-tab":
			API.FocusTab(query)
			return
		case "new-window":
			API.OpenNewWindow()
			return
		case "new-tab":
			API.OpenNewTab()
			return

		// Change setting
		case "cache-favicons":
			API.CacheFavicons()
			return
		case "select-profile":
			API.SelectProfile(wf, query)
		case "change-profile":
			API.ChangeProfile(query)
			return
		case "select-browser":
			API.SelectBrowser(wf, query)
		case "change-browser":
			API.ChangeBrowser(query)
			return
		}

		wf.SendFeedback()
	} else {
		panic("Improper arguments")
	}
}

func init() {
	repoUrl := "jopemachine/alfred-chromium-workflow"
	wf = aw.New(update.GitHub(repoUrl), aw.HelpURL(repoUrl + "/issues"), aw.MaxResults(200))
}

func main() {
	API.ImportConfig()
	lctime.SetLocale(API.Conf.Locale)
	wf.Run(alfredCallback)
}
