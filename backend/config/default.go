package config

var defaultCFG = CFG{
	VersionCode:   "Testify",
	AppVersion:    "v0.0.2",
	HlaeVersion:   "",
	FFmpegVersion: "",
	HlaeAPI:       "https://api.github.com/repos/advancedfx/advancedfx/releases/latest",
	HlaeCdnAPI:    "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master/api.json",
	FFmpegAPI:     "https://www.gyan.dev/ffmpeg/builds", // builds/ffmpeg-release-essentials.7z
	FFmpegCdnAPI:  "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/api.json",
	HlaePath:      "",
	Init:          false,
	Standalone:    false,
	HlaeState:     false,
	FFmpegState:   false,
}
