<h1 align="center">
  <a href="https://www.npmjs.com/package/alfred-chromium-workflow">
		Alfred-chromium-workflow<br>
	<img src="https://img.shields.io/badge/Alfred-4-blueviolet">
	<img src="https://img.shields.io/github/downloads/jopemachine/alfred-chromium-workflow/total.svg">
	<img src="https://img.shields.io/github/license/jopemachine/alfred-chromium-workflow.svg" alt="License">
  </a>
</h1>

Alfred workflow for Chromium browsers

## Why?

* Support favicon images in almost all features

* Support switching browsers, profiles

* Lightning-fast

* Provide localized subtitle

## 🌈 Features

* 📄 *Search Visit History `(chh)`*

![](./imgs/chh.png)

* 🔖 *Search Bookmark, bookmark folders `(chb, chf)`*

![](./imgs/chb.png)

* 📁 *Search Bookmark folder `(chf)`*

![](./imgs/chf.png)

* 📜 *Search Search query history `(chs)`*

![](./imgs/chs.png)

* 🔎 *Search Download logs `(chd)`*

![](./imgs/chd.png)

* 📒 *Search and Copy Your Autofill data `(cha)`*

![](./imgs/cha.png)

* *Search Your Opened Tabs and Focus or Close Them `(cht)`*

## 📌 Prerequisite

The prerequisites below are required to use that package.

* [Alfred powerpack](https://www.alfredapp.com/powerpack/)

##  🔨 How to install

Download and open `alfredworkflow` file in the [Release page]().

## 📗 Usage

### chb

Search bookmark.

### chf

Search bookmark folders.

### chh

Search visit history.

You can append `#` to search word to search only the logs in that `url`.

Example:

`chh #youtube [some_word_to_search]`

### chd

Search download history.

### chs

Search your query based on visit history.

You can append `#` to search word to search only the logs in that `url`.

Example:

`chs #github [some_word_to_search]`

### cha

Search chrome autofill data.

### chid

Search chrome's login data (including email).

### ch > browser 

You can switch browsers with this command with ease

### ch > profile

You can switch profiles with this command with ease

### ch > cache

Cache favicon images from your visit history in advance.

### ch > clearcache

Clear favicon cache.

## 🔖 Options

You can configure below options through Workflow Environment Variables.

### Profile

Type: `string`

Your browser's profile name.

### SwitchableProfiles

Type: `string`

All switchable profiles.

Each profile name should be splited with comma (`,`).

### Locale

Type: `string (enum)`

Possible values: Refer to the below page for seeing which locales are supported.

https://github.com/klauspost/lctime/tree/master/internal/locales

### Browser

Type: `string (enum)`

Select the browser to which you want the workflow to the workflow.

Possible values: `Chrome`, `Chrome Canary`, `Edge`, `Chromium`, `Brave`

### ExcludeDomains

Type: `string[]`

You can exclude specific domains in your result. 

This is not applied to `chb`.

### ResultLimit

Type: `number`

Displays as many search results.

## License

MIT © [jopemachine](https://github.com/jopemachine/alfred-chromium-workflow)
