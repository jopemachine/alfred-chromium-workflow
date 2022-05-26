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
		case "download":
			API.FetchDownloadHistory(wf, query)
		case "login":
		case "autofill":
			API.FetchAutofillData(wf, query)
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

