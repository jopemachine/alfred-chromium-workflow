package src

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deanishe/awgo"
)

var FetchBookmarkFolder = func(wf *aw.Workflow, query string) {
	InitBookmarkJsonTraversal()
	bookmarkRoot := GetChromeBookmark()

	bookmarkFolders := TraverseBookmarkJSONObject(bookmarkRoot, TraverseBookmarkJsonOption{Targets: []string{"folder"}, Depth: 99})
	sort.Slice(bookmarkFolders, func(i, j int) bool {
		if strings.Compare(bookmarkFolders[i].Name, bookmarkFolders[j].Name) == 1 {
			return true
		} else {
			return false
		}
	})

	for _, folder := range bookmarkFolders {
		folderChildLen := 0
		if folder.Children != nil {
			for _, child := range folder.Children {
				if child.(map[string]interface{})["type"] == "url" {
					folderChildLen += 1
				}
			}
		}

		wf.NewItem(folder.Name).
			Valid(true).
			Subtitle(fmt.Sprintf(`%d items included`, folderChildLen)).
			Arg(folder.Id).
			Autocomplete(folder.Name).
			Copytext(folder.Name).
			Largetype(folder.Name).
			Icon(&aw.Icon{"assets/folder.png", ""}).
			Var("folder", fmt.Sprintf(`--%s=%s`, "folderId", folder.Id))
	}

	if query != "" {
		wf.Filter(query)
	}
}
