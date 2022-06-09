package src

import (
	"fmt"
	"github.com/deanishe/awgo"
	"sort"
)

var FetchBookmark = func(wf *aw.Workflow, query string) {
	InitBookmarkJsonTraversal()
	bookmarkRoot := GetChromeBookmark()
	input, flags := ParseQueryFlags(query)
	var bookmarks []BookmarkItem

	if folderId, ok := flags["folderId"]; ok {
		folders := TraverseBookmarkJSONObject(bookmarkRoot, TraverseBookmarkJsonOption{Targets: []string{"folder"}, Depth: 99})

		for _, folder := range folders {
			if folder.Id == folderId {
				bookmarks = TraverseBookmarkArray(folder.Children, TraverseBookmarkJsonOption{Targets: []string{"url"}, Depth: 1})
			}
		}
		if bookmarks == nil {
			panic(fmt.Sprintf("folderId not found: %s", folderId))
		}
	} else {
		bookmarks = TraverseBookmarkJSONObject(bookmarkRoot, TraverseBookmarkJsonOption{Targets: []string{"url"}, Depth: 99})
	}

	historyDB := GetHistoryDB(wf)
	visitHistories, err := historyDB.Query("SELECT url FROM urls")
	CheckError(err)

	visitFrequency := make(map[string]int)

	for visitHistories.Next() {
		var url string
		err := visitHistories.Scan(&url)
		CheckError(err)

		visitFrequency[url] += 1
	}

	sort.Slice(bookmarks, func(i, j int) bool {
		ithFreq := visitFrequency[bookmarks[i].Url]
		jthFreq := visitFrequency[bookmarks[j].Url]

		if ithFreq > 0 && jthFreq > 0 {
			if ithFreq > jthFreq {
				return true
			} else {
				return false
			}
		}

		if ithFreq > 0 {
			return true
		}

		return false
	})

	for _, bookmark := range bookmarks {
		domainName := ExtractDomainName(bookmark.Url)
		iconPath := fmt.Sprintf(`%s/%s.png`, GetFaviconDirectoryPath(wf), domainName)
		CheckError(err)

		item := wf.NewItem(bookmark.Name).
			Valid(true).
			Subtitle(bookmark.Url).
			Quicklook(bookmark.Url).
			Arg(bookmark.Url).
			Copytext(bookmark.Url).
			Autocomplete(bookmark.Name).
			Largetype(bookmark.Name)

		item.Cmd().Subtitle("Press Enter to copy this url to clipboard")

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}
	}

	if input != "" {
		wf.Filter(input)
	}
}
