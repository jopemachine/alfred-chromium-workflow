package src

// https://github.com/nikitavoloboev/alfred-web-searches/blob/master/update.go

import (
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
)

func doUpdate() error {
	log.Println("Checking for update...")
	return wf.CheckForUpdate()
}

// checkForUpdate runs "./alsf update" in the background if an update check is due.
func checkForUpdate() error {
	if !wf.UpdateCheckDue() || wf.IsRunning("update") {
		return nil
	}
	cmd := exec.Command(os.Args[0], "update")
	return wf.RunInBackground("update", cmd)
}

// showUpdateStatus adds an "update available!" message to Script Filters if an update is available
// and query is empty.
func ShowUpdateStatus(wf *aw.Workflow, query string) {
	if query != "" {
		return
	}

	if wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		log.Println("Workflow update available!")

		wf.NewItem("Workflow update available!").
			Subtitle("⇥ or ↩ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(&aw.Icon{Value: "icons/update-available.png"})
	}
}
