package src

// Ref: https://github.com/deanishe/alfred-sublime-text/blob/master/cli.go

import (
	"log"

	aw "github.com/deanishe/awgo"
)

var (
	iconForum           = &aw.Icon{Value: "assets/forum.png"}
	iconHelp            = &aw.Icon{Value: "assets/help.png"}
	iconIssue           = &aw.Icon{Value: "assets/issue.png"}
	iconUpdateAvailable = &aw.Icon{Value: "assets/update-available.png"}
	iconUpdateOK        = &aw.Icon{Value: "assets/update-ok.png"}
)

var UpdateWorkflow = func(wf *aw.Workflow) {
	if wf.UpdateAvailable() {
		wf.InstallUpdate()
	}
}

var RunWorkflowHelper = func(wf *aw.Workflow, query string) {
	wf.CheckForUpdate()
	wf.Configure(aw.SuppressUIDs(true))

	if wf.UpdateAvailable() {
		log.Println("Workflow Update Available!")

		wf.NewItem("Workflow Update Available!").
			Subtitle("⇥ or ↩ to install update").
			Valid(true).
			Icon(iconUpdateAvailable).
			Var("update", "true")
	} else {
		wf.NewItem("Workflow Is Up To Date").
			Subtitle("").
			Valid(false).
			Icon(iconUpdateOK)
	}

	wf.NewItem("View Help File").
		Subtitle("Open workflow help in your browser").
		Valid(true).
		Icon(iconHelp).
		Var("url", "https://github.com/jopemachine/alfred-chromium-workflow/blob/master/README.md")

	wf.NewItem("Report Issue").
		Subtitle("Open workflow issue tracker in your browser").
		Icon(iconIssue).
		Var("url", "https://github.com/jopemachine/alfred-chromium-workflow/issues").
		Valid(true)

	wf.NewItem("Visit Forum Thread").
		Subtitle("Open workflow thread on alfredforum.com in your browser").
		Icon(iconForum).
		Var("url", "https://www.alfredforum.com/topic/18380-chromium-based-browser-workflow-supporting-browser-profile-switching/").
		Valid(true)

	wf.NewItem("Cache Favicons Manually").
		Subtitle("").
		Valid(true).
		Var("cacheFavicons", "true")

	if query != "" {
		wf.Filter(query)
	}
}
