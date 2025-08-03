package core

import (
	"context"
	"fmt"
	"os/exec"
)

// VideoProcessor định nghĩa interface xử lý video
// Triển khai bằng ffmpeg qua shell

type VideoProcessor interface {
	TranscodeToHLS(ctx context.Context, inputPath, outputDir string, qualities []string) error
}

// FFMPEGVideoProcessor là implement VideoProcessor dùng ffmpeg

type FFMPEGVideoProcessor struct{}

func NewFFMPEGVideoProcessor() *FFMPEGVideoProcessor {
	return &FFMPEGVideoProcessor{}
}

// TranscodeToHLS chuyển video sang HLS với nhiều chất lượng
func (p *FFMPEGVideoProcessor) TranscodeToHLS(ctx context.Context, inputPath, outputDir string, qualities []string) error {
	// qualities ví dụ: ["360p", "480p", "720p"]
	// mapping chất lượng sang thông số ffmpeg
	qualityMap := map[string]string{
		"360p":  "-vf scale=-2:360 -b:v 800k",
		"480p":  "-vf scale=-2:480 -b:v 1400k",
		"720p":  "-vf scale=-2:720 -b:v 2800k",
		"1080p": "-vf scale=-2:1080 -b:v 5000k",
	}
	for _, q := range qualities {
		ffmpegArgs := []string{"-i", inputPath}
		if opt, ok := qualityMap[q]; ok {
			ffmpegArgs = append(ffmpegArgs, splitArgs(opt)...) // scale, bitrate
		}
		ffmpegArgs = append(ffmpegArgs,
			"-c:a", "aac", "-ar", "48000", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0",
			"-g", "48", "-keyint_min", "48", "-hls_time", "4", "-hls_playlist_type", "vod",
			"-f", "hls",
			fmt.Sprintf("%s/%s.m3u8", outputDir, q),
		)
		cmd := exec.CommandContext(ctx, "ffmpeg", ffmpegArgs...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("ffmpeg error (%s): %v", out, err)
		}
	}
	return nil
}

// splitArgs tách chuỗi thành slice cho exec.Command
func splitArgs(s string) []string {
	var out []string
	for _, v := range []rune(s) {
		if v == ' ' {
			continue
		}
		out = append(out, string(v))
	}
	return []string{}
}
