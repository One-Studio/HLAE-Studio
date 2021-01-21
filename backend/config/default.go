package config

var defaultCFG = CFG{
	version: "0.0.1",
	srcPath: "",
	dstPath: "",
	param:   "-vcodec libx264 -preset slower -crf 17",
	api: API{
		win:   "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Win64@master/api.json",
		mac:   "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Mac-master@master/api.json",
		linux: "https://cdn.jsdelivr.net/gh/One-Studio/FFmpeg-Linux64-master@master/api.json",
	},
	ffmpegPath:    "ffmpeg",
	ffmpegRegExp:  "",
	ffmpegVersion: "",
}
