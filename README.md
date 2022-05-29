<h1 align="center">
		Alfred-chromium-workflow<br>
		<img src="https://img.shields.io/badge/Alfred-4-blueviolet">
		<img src="https://img.shields.io/github/downloads/jopemachine/alfred-chromium-workflow/total.svg">
		<img src="https://img.shields.io/github/license/jopemachine/alfred-chromium-workflow.svg" alt="License">
</h1>

Alfred workflow for Chromium browsers

## Why?

This workflow originated from [alfred-chrome-workflow](https://github.com/jopemachine/alfred-chrome-workflow).

Previous workflow has several [npm related installation issues like this](https://github.com/jopemachine/alfred-chrome-workflow/issues/13#issuecomment-1103938917).

In addition to resolving these issues, this workflow has following several advantages over the previous one.

* Support favicon images in almost all features
* Support switching browsers, profiles with ease
* Support more Chromium based browsers
* Lightning-fast
* Provide localized subtitle

## 🌈 Features

* 📄 *Search Visit History*
* 🔖 *Search Bookmark, bookmark folders*
* 📁 *Search Bookmark folder*
* 📜 *Search Search query history*
* 🔎 *Search Download logs*
* 📒 *Search and Copy Your Autofill data*
* 📎 *Search Your Opened Tabs and Focus or Close Them*

## 📌 Prerequisite

The prerequisites below are required to use that package.

* [Alfred Powerpack](https://www.alfredapp.com/powerpack/)

## 🎯 Supported Browsers

Currently supported browsers are as follows.

* `Chrome`
* `Chrome Canary`
* `Chromium`
* `Edge`
* `Brave`
* `Naver Whale`
* `Epic`

## 🔨 How to install

Download and open `alfredworkflow` file in the [Release page]().

## 📗 Usage

### chb

Retrieve bookmarks.

![](./imgs/chb.png)

### chf

Retrieve bookmark folders.

![](./imgs/chf.png)

### chh

Retrieve visit histories.

Append `#` to search word to search only the logs in that `url`.
 
Example:

`chh #youtube [some_word_to_search]`

![](./imgs/chh.png)

### chd

Retrieve download histories.

![](./imgs/chd.png)

### chdc

Retrieve download histories but only shows existing files.

### chs

Retrieve your search query based on visit histories.

Append `#` to search word to search only the logs in that `url`.

Example:

`chs #github [some_word_to_search]`

![](./imgs/chs.png)

### cha

Retrieve autofill data.

### chid

Retrieve login data (including email).

### cho

Open new window through selected profiles.

If you change your profile, other commands try to work with previous profile.

This command would be useful in such case.

Open new window with changed profile before the command.

### cht

Search opened tabs and focus, close them.

### ch browser

Switch browsers with ease.

This command also requires you to change the browser profile.

### ch profile

Switch profiles with ease.

### ch cache

Cache favicon images from your visit history in advance.

You probably not need this because caching runs in background automatically.

## 🔖 Options

Configure below options through Workflow Environment Variables.

![](./imgs/conf.png)

### Profile

Browser's profile name.

You can change this value through `ch profile` with ease.

### SwitchableProfiles

List up your custom profile names here.

You can switch your profile through `ch profile` with `Profile {number}` and these values.

Each profile name should be splited with comma (`,`).

### Locale

Refer to the following page for seeing supported locales.

https://github.com/klauspost/lctime/tree/master/internal/locales

### Browser

Browser name.

You can change this value through `ch browser` with ease.

### ResultCountLimit

Max number of items to show in Alfred.

## 🌟 Contribution

Contributions of any kind are welcome.

### Add supporting new browser

This workflow needs your help to support as many browsers as possible.

To add new browser, follow below guideline.

1. Check the browser is based on Chromium

2. Add new path to `GetProfileRootPath` in `src/utils.go`.

3. Add the browser's proper Application Name to `getApplicationName` in `src/tabManager.go`. You can check this value through Applescript Editor's Dictionary.

4. Add proper Open URL block of `info.plist` using Alfred.

5. Add new browser's item `SelectBrowser` in `src/config.go`.

6. Add the browser name to `READMD.md`

## Related

You may also consider below workflows interesting.

- [chrome-control](https://github.com/bit2pixel/chrome-control): A JXA script and an Alfred Workflow for controlling Google Chrome

## License

MIT © [jopemachine](https://github.com/jopemachine/alfred-chromium-workflow)
