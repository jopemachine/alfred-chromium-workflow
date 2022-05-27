package main

import (
	"flag"

	"github.com/deanishe/awgo"
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
		case "cache-favicons":
			API.CacheFavicons()
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
	wf = aw.New()
}

func main() {
	API.ImportConfig()
	lctime.SetLocale(API.Conf.Locale)
	wf.Run(alfredCallback)
}
