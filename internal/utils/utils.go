package utils

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	regexPattern = `https:\/\/open\.spotify\.com\/playlist\/([a-zA-Z0-9]+)`
)

func ExtractPlaylistID(link string) string {
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(link)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func ExtractYouTubeLink(input string) string {
	re := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([^&]+)`)
	match := re.FindStringSubmatch(input)
	if len(match) < 2 {
		return "" // No match found
	}
	videoID := match[1]
	return "https://www.youtube.com/watch?v=" + videoID
}

func DownloadMusic(videoURL string, downloadComplete chan<- string) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Specify the download folder
	downloadFolder := filepath.Join(cwd, "downloads")

	// Ensure download folder exists, create if not
	err = os.MkdirAll(downloadFolder, 0755)
	if err != nil {
		return err
	}

	// Command to get the video title
	cmd := exec.Command("youtube-dl", "--get-title", "--skip-download", "--", videoURL)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Extract the title from the output
	title := strings.TrimSpace(out.String())
	if title == "" {
		title = "Untitled" // Default title if not found
	}

	// Command to download music with metadata and thumbnail using youtube-dl
	cmd = exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "--embed-thumbnail", "--add-metadata", "-o", filepath.Join(downloadFolder, "%(title)s.%(ext)s"), videoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the download
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Get the downloaded file name
	downloadedFileName := filepath.Join(downloadFolder, title+".mp3")

	downloadComplete <- downloadedFileName
	return nil
}
