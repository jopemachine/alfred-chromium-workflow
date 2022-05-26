package src

import (
	"fmt"

	"github.com/deanishe/awgo"
)

var FetchBookmarkFolder = func (wf *aw.Workflow, query string) {
	InitBookmarkJsonTraversal()
	bookmarkRoot := GetChromeBookmark()

	bookmarkFolders := TraverseBookmarkJSONObject(bookmarkRoot, TraverseBookmarkJsonOption{ Targets: []string{"folder"}, Depth: 99 })

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
			Icon(&aw.Icon{"assets/folder.png", ""}).
			Var("folder", fmt.Sprintf(`--%s=%s`, "folderId", folder.Id))
	}

	wf.Filter(query)
}
