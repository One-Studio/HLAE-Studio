package config

var defaultCFG = CFG{
	VersionCode:   "Testify",
	AppVersion:    "v0.0.1",
	HlaeVersion:   "",
	FFmpegVersion: "",
	HlaeAPI:       "https://api.github.com/repos/advancedfx/advancedfx/releases/latest",
	HlaeCdnAPI:    "https://cdn.jsdelivr.net/gh/One-Studio/HLAE-Archive@master",
	FFmpegAPI:     "https://www.gyan.dev/ffmpeg/builds/release-version", // builds/ffmpeg-release-essentials.7z
	FFmpegCdnAPI:  "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master",
	HlaePath:      "",
	Init:          false,
	HlaeState:     false,
	FFmpegState:   false,
}
