package main

import (
	"flag"

	"github.com/klauspost/lctime"
	"github.com/deanishe/awgo"
	API "github.com/jopemachine/alfred-chromium-workflow/src"
)

var wf *aw.Workflow

func alfredCallback () {
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		commandType, query := args[0], args[1]

		switch commandType {
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

		// Tab related features
		case "listup-tabs":
			API.ListOpenedTabs(query)
			return
		case "close-tab":
			API.CloseTab(query)
			return
		case "focus-tab":
			API.FocusTab(query)
			return

		// Setting
		case "select-browser":
		}

		wf.SendFeedback()
	} else {
	}
}

func init() {
	wf = aw.New()
}

func main () {
	API.ImportConfig()
	lctime.SetLocale(API.Conf.Locale)
	wf.Run(alfredCallback)
}

