package config

var defaultCFG = CFG{
	VersionCode:   "Testify",
	AppVersion:    "v1.1.0",
	HlaeVersion:   "",
	FFmpegVersion: "",
	HlaeAPI:       "https://api.upup.cool/get/hlae",
	//HlaeCdnAPI:    "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/api.json",
	FFmpegAPI:     "https://api.upup.cool/get/ffmpeg", // builds/ffmpeg-release-essentials.7z
	//FFmpegCdnAPI:  "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/api.json",
	HlaePath:      "",
	Init:          false,
	Standalone:    false,
	HlaeState:     false,
	FFmpegState:   false,
}
