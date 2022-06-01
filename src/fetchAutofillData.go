package src

import (
	"fmt"
	"github.com/deanishe/awgo"
)

var FetchAutofillData = func(wf *aw.Workflow, query string) {
	var dbQuery = fmt.Sprintf(`
		SELECT value, name, date_created, count
			FROM autofill
			WHERE value LIKE '%%%s%%' OR name LIKE '%%%s%%'
			ORDER BY count DESC
	`, query, query)

	webDataDB := GetWebDataDB(wf)
	rows, err := webDataDB.Query(dbQuery)
	CheckError(err)

	var autofillValue string
	var autofillLabel string
	var createdDate int64
	var count int

	for rows.Next() {
		err := rows.Scan(&autofillValue, &autofillLabel, &createdDate, &count)
		CheckError(err)

		var subtitle string

		unixTimestamp := ConvertChromeTimeToUnixTimestamp(createdDate)
		localeTimeStr := GetLocaleString(unixTimestamp)

		subtitle += fmt.Sprintf(`Label: '%s', Created in '%s'`, autofillLabel, localeTimeStr)

		item := wf.NewItem(autofillValue).
			Subtitle(subtitle).
			Valid(true).
			Copytext(autofillValue).
			Arg(autofillValue).
			Largetype(autofillValue).
			Autocomplete(autofillValue).
			Icon(&aw.Icon{"assets/info.png", ""})

		item.Cmd().Subtitle("Press Enter to paste this value directly")
	}
}
