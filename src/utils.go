package src

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/deanishe/awgo"
	"github.com/klauspost/lctime"
	_ "github.com/mattn/go-sqlite3"
	psl "github.com/weppos/publicsuffix-go/publicsuffix"
)

var CheckError = func(err error) {
	if err != nil {
		panic(err)
	}
}

var GetProfileRootPath = func(browserName string) string {
	var targetPath string

	user, err := user.Current()
	CheckError(err)
	userName := user.Username

	switch browserName {
	case "Opera":
		targetPath = `/Users/%s/Library/Application Support/com.operasoftware.Opera`
	case "Chrome Canary":
		targetPath = `/Users/%s/Library/Application Support/Google/Chrome Canary`
	case "Edge":
		targetPath = `/Users/%s/Library/Application Support/Microsoft Edge`
	case "Edge Canary":
		targetPath = `/Users/%s/Library/Application Support/Microsoft Edge Canary`
	case "Chromium":
		targetPath = `/Users/%s/Library/Application Support/Google/Chrome Cloud Enrollment`
	case "Brave":
		targetPath = `/Users/%s/Library/Application Support/BraveSoftware/Brave-Browser`
	case "Chrome":
		targetPath = `/Users/%s/Library/Application Support/Google/Chrome`
	case "Chrome Beta":
		targetPath = `/Users/%s/Library/Application Support/Google/Chrome Beta`
	case "Naver Whale":
		targetPath = `/Users/%s/Library/Application Support/Naver/Whale`
	case "Vivaldi":
		targetPath = `/Users/%s/Library/Application Support/Vivaldi`
	case "Epic":
		targetPath = `/Users/%s/Library/Application Support/HiddenReflex/Epic`
	default:
		panic("Unsupported browser. Please consider to make a issue to support the browser if the browser is based on Chromium.")
	}

	return fmt.Sprintf(targetPath, userName)
}

var GetDBFilePath = func(browserName string, chromeProfilePath string, dbFile string) string {
	if browserName == "Opera" {
		return fmt.Sprintf(`%s/%s`, GetProfileRootPath(browserName), dbFile)
	}
	return fmt.Sprintf(`%s/%s/%s`, GetProfileRootPath(browserName), chromeProfilePath, dbFile)
}

// Used in `chs`, `chh`
var ParseUserQuery = func(query string) (titleQuery string, domainQuery string, isDomainSearch bool) {
	titleQuery = ""
	domainQuery = ""
	isDomainSearch = false

	// Useless since `chm` not implemented
	// artistQuery = ""
	// isArtistSearch = false

	if strings.Contains(query, "#") || strings.Contains(query, "@") {
		var words = strings.Split(query, " ")

		for _, word := range words {
			// else if strings.HasPrefix(word, "@") && len(word) > 1 {
			// 	isArtistSearch = true
			// 	artistQuery = word[1 : len(word)-1]
			// }

			if strings.HasPrefix(word, "#") && len(word) > 1 {
				isDomainSearch = true
				domainQuery = word[1 : len(word)-1]
			} else {
				// TODO: Refactor below logic using `strings.Join`
				if titleQuery == "" {
					titleQuery += word
				} else {
					titleQuery += " " + word
				}
			}
		}
	} else {
		titleQuery = query
	}

	return
}

func CopyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		panic(err)
	}
	err = out.Sync()
	if err != nil {
		panic(err)
	}
}

var GetHistoryDB = func(wf *aw.Workflow) *sql.DB {
	var targetPath = GetDBFilePath(Conf.Browser, Conf.Profile, "History")
	dest := filepath.Join(GetTempDataPath(wf), CONSTANT.HISTORY_DB)
	CopyFile(targetPath, dest)
	db, err := sql.Open("sqlite3", dest)
	CheckError(err)

	return db
}

var GetFaviconDB = func(wf *aw.Workflow) *sql.DB {
	var targetPath = GetDBFilePath(Conf.Browser, Conf.Profile, "Favicons")
	dest := filepath.Join(GetTempDataPath(wf), CONSTANT.FAVICON_DB)
	CopyFile(targetPath, dest)
	db, err := sql.Open("sqlite3", dest)
	CheckError(err)

	return db
}

var GetWebDataDB = func(wf *aw.Workflow) *sql.DB {
	var targetPath = GetDBFilePath(Conf.Browser, Conf.Profile, "Web Data")
	dest := filepath.Join(GetTempDataPath(wf), CONSTANT.WEB_DATA_DB)
	CopyFile(targetPath, dest)
	db, err := sql.Open("sqlite3", dest)
	CheckError(err)

	return db
}

var GetLoginDataDB = func(wf *aw.Workflow) *sql.DB {
	var targetPath = GetDBFilePath(Conf.Browser, Conf.Profile, "Login Data")
	dest := filepath.Join(GetTempDataPath(wf), CONSTANT.WEB_DATA_DB)
	CopyFile(targetPath, dest)
	db, err := sql.Open("sqlite3", dest)
	CheckError(err)

	return db
}

var GetChromeBookmark = func() map[string]interface{} {
	var bookmarkJson map[string]interface{}
	var bookmarkFilePath = GetDBFilePath(Conf.Browser, Conf.Profile, "Bookmarks")

	bookmarkData, err := ioutil.ReadFile(bookmarkFilePath)
	CheckError(err)
	err = json.Unmarshal(bookmarkData, &bookmarkJson)
	CheckError(err)

	return bookmarkJson["roots"].(map[string]interface{})
}

// Ref: https://golangcode.com/how-to-check-if-a-string-is-a-url/
// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

var ExtractDomainName = func(url string) (domainName string) {
	if !isValidUrl(url) {
		return "unknown"
	}

	var hostname string
	if strings.Contains(url, "//") {
		hostname = strings.Split(url, "/")[2]
	} else {
		hostname = strings.Split(url, "/")[0]
	}

	hostname = strings.Split(hostname, ":")[0]
	hostname = strings.Split(hostname, "?")[0]

	domainName, err := psl.Domain(hostname)
	if err != nil {
		return hostname
	}

	return domainName
}

var ConvertChromeTimeToUnixTimestamp = func(time int64) int64 {
	return int64((math.Floor((float64(time) / 1000000) - 11644473600)))
}

var FileExist = func(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

var GetLocaleString = func(unixTime int64) string {
	return lctime.Strftime("%c", time.Unix(unixTime, 0))
}

// Used only in fetchBookmark.go of`chf` command
var ParseQueryFlags = func(userQuery string) (input string, options map[string]string) {
	options = make(map[string]string)

	for _, args := range strings.Split(userQuery, " ") {
		reg, err := regexp.Compile("--[a-zA-Z\\d]*")
		CheckError(err)
		argList := strings.Split(args, " ")

		for _, arg := range argList {
			if reg.MatchString(arg) {
				key := strings.Split(strings.Split(arg, "--")[1], "=")[0]
				value := strings.Split(arg, "=")[1]
				options[key] = value
			} else {
				input += (arg + " ")
			}
		}
	}

	if strings.HasSuffix(input, " ") {
		input = strings.Trim(input, " ")
	}

	return
}

// TODO: Replace below function with stdlib's one when it is merged
// Ref: https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func StringContains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

type BookmarkItem struct {
	Id       string        `json:"id"`
	Url      string        `json:"url,omitempty"`
	Name     string        `json:"name,omitempty"`
	Children []interface{} `json:"children,omitempty"`
}

type TraverseBookmarkJsonOption struct {
	Targets []string
	Depth   int
}

// InitBookmarkJsonTraversal should be called first before calling this function
var TraverseBookmarkJSONObject func(item map[string]interface{}, options TraverseBookmarkJsonOption) []BookmarkItem
var TraverseBookmarkArray func(item []interface{}, options TraverseBookmarkJsonOption) []BookmarkItem

var InitBookmarkJsonTraversal = func() {
	TraverseBookmarkJSONObject = func(jsonObject map[string]interface{}, options TraverseBookmarkJsonOption) (result []BookmarkItem) {
		// Base case
		if options.Depth <= -1 {
			return []BookmarkItem{}
		}

		// Base case
		if jsonObject["type"] == "url" {
			if StringContains(options.Targets, "url") {
				return []BookmarkItem{
					{
						jsonObject["id"].(string),
						jsonObject["url"].(string),
						jsonObject["name"].(string),
						nil,
					},
				}
			}

			return []BookmarkItem{}
		}

		if StringContains(options.Targets, "folder") && jsonObject["type"] == "folder" {
			result = append(result, BookmarkItem{
				jsonObject["id"].(string),
				"",
				jsonObject["name"].(string),
				jsonObject["children"].([]interface{}),
			})

			childResult := TraverseBookmarkArray(jsonObject["children"].([]interface{}), options)
			result = append(result, childResult...)
			return result
		}

		target := jsonObject

		for _, child := range target {
			switch child.(type) {
			case map[string]interface{}:
				childResult := TraverseBookmarkJSONObject(child.(map[string]interface{}), options)
				result = append(result, childResult...)
			case []interface{}:
				childResult := TraverseBookmarkArray(child.([]interface{}), options)
				result = append(result, childResult...)
			}
		}

		return result
	}

	TraverseBookmarkArray = func(item []interface{}, options TraverseBookmarkJsonOption) []BookmarkItem {
		// Base case
		if options.Depth <= -1 {
			return []BookmarkItem{}
		}

		target := item
		result := []BookmarkItem{}

		for _, child := range target {
			switch child.(type) {
			case map[string]interface{}:
				childResult := TraverseBookmarkJSONObject(child.(map[string]interface{}), options)
				result = append(result, childResult...)
			case []interface{}:
				childResult := TraverseBookmarkArray(child.([]interface{}), options)
				result = append(result, childResult...)
			}
		}

		return result
	}
}

var CheckBrowserIsInstalled = func(browserName string) bool {
	return FileExist(GetProfileRootPath(browserName))
}

// Ref: https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func IsEmptyDirectory(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
