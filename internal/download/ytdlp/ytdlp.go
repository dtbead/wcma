package ytdlp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

var regexpYoutubeURL *regexp.Regexp

func init() {
	regexpYoutubeURL = regexp.MustCompile(`youtube\.com\/watch\?v=([0-9A-Za-z_-]{10}[048AEIMQUYcgkosw])`)
}

func isValidURL(url string) bool {
	return regexpYoutubeURL.MatchString(url)
}

func cleanURL(url string) string {
	return "https://youtube.com/watch?v=" + regexpYoutubeURL.FindStringSubmatch(url)[1]
}

type Ytdlp struct {
	cookies []byte
}

func NewYtdlp(netscape_cookies []byte) Ytdlp {
	return Ytdlp{cookies: netscape_cookies}
}

func (y Ytdlp) Download(ctx context.Context, url string, output io.Writer) (youtube *entities.Youtube, extension string, err error) {
	if !isValidURL(url) {
		return nil, "", errors.New("invalid youtube url")
	}
	url = cleanURL(url)

	currentTime := fmt.Sprint(time.Now().UnixMilli())
	file_output := os.TempDir() + "\\" + currentTime + `_%(id)s.%(ext)s`

	args := []string{
		"--ignore-config",
		"--no-playlist",
		"--output",
		file_output,
		"--restrict-filenames",
		"--no-simulate",
		"--no-part",
		"--no-write-comments",
		"--no-cache-dir",
		"--no-write-thumbnail",
		"--no-embed-metadata",
		"--no-embed-info-json",
		//"-S",
		// "res:480",
	}
	args = append(args, "-J", "--print", file_output, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, "", errors.Join(err, errors.New(stderr.String()))
	}

	// yt-dlp will print the video file path with a trailing newline
	file_output = strings.Split(stdout.String(), "\n")[0]
	defer os.Remove(file_output)

	// then it will print the json metadata of said video afterwards
	json := strings.Split(stdout.String(), "\n")[1]

	m, err := newMetadata([]byte(json))
	if err != nil {
		return nil, "", errors.Join(err, errors.New(stderr.String()))
	}

	f, err := os.Open(file_output)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	_, err = io.Copy(output, f)
	if err != nil {
		return nil, "", err
	}

	yt := m.ToYoutubeEntity()
	ext := m.Extension

	stdout.Reset()
	stderr.Reset()
	m = nil

	return &yt, ext, nil
}
