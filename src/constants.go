package src

type constant struct {
	HISTORY_DB       string
	FAVICON_DB       string
	MEDIA_HISTORY_DB string
	WEB_DATA_DB      string
	LOGIN_DATA_DB    string
	COOKIE_DB        string
}

var CONSTANT = constant{
	HISTORY_DB:       "_history.db",
	FAVICON_DB:       "_favicon.db",
	MEDIA_HISTORY_DB: "_mediaHistory.db",
	WEB_DATA_DB:      "_webData.db",
	LOGIN_DATA_DB:    "_loginData.db",
	COOKIE_DB:        "_cookies.db",
}
