package api

import "encoding/xml"

type ReleaseDelivr struct {
	Version      string   `json:"Version"`
	VersionList  []string `json:"VersionList"`
	ReleaseTime  string   `json:"ReleaseTime"`
	CheckTime    string   `json:"CheckTime"`
	DownloadLink []string `json:"DownloadLink"`
	Format       int      `json:"Format"`
	ReleaseNote  string   `json:"ReleaseNote"`
}

//Github Asset
type Asset struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

//Github latest info
type GitHubLatest struct {
	URL     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Assets  []Asset `json:"assets"`
}

//HLAE changelog.xml
type Changelog struct {
	XMLName xml.Name `xml:"changelog"`
	Text    string   `xml:",chardata"`
	Release []struct {
		Text    string `xml:",chardata"`
		Name    string `xml:"name"`
		Version string `xml:"version"`
		Time    string `xml:"time"`
		Changes struct {
			Text   string `xml:",chardata"`
			Change []struct {
				Text string   `xml:",chardata"`
				Type string   `xml:"type,attr"`
				Br   []string `xml:"br"`
			} `xml:"change"`
			Changed struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"changed"`
		} `xml:"changes"`
		Comments struct {
			Text string   `xml:",chardata"`
			Br   []string `xml:"br"`
		} `xml:"comments"`
	} `xml:"release"`
	H1 string `xml:"h1"`
}

//Github FFmpeg info
type FFmpegTag struct {
	Message    string `json:"message"`
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}
