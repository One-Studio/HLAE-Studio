package api

type OsMirror struct {
	Version      string   `json:"Version"`
	VersionList  []string `json:"VersionList"`
	ReleaseTime  string   `json:"ReleaseTime"`
	CheckTime    string   `json:"CheckTime"`
	DownloadLink []string `json:"DownloadLink"`
	Format       int      `json:"Format"`
	ReleaseNote  string   `json:"ReleaseNote"`
}
