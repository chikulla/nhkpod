package nhkpod

import (
	"errors"
	"github.com/bogem/id3v2"
	"github.com/canhlinh/hlsdl"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Episode struct {
	ID          string
	M3u8Url     string
	Started     time.Time
	Program     string
	Title       string
	Description string
	FilePath    string
	Size        int64
}

var TFormat = "200601021504"
var TFormatZone = TFormat + "-0700"
var Encoding = ".aac"

func NewFromFile(filePath string, size int64) (*Episode, error) {
	f, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	defer f.Close()
	if err != nil {
		return nil, err
	}
	id, t, err := FileNameToIDAndStarted(filePath)
	return &Episode{
		ID:          *id,
		Started:     *t,
		Title:       f.Title(),
		Program:     f.Artist(),
		Description: f.Album(),
		FilePath:    filePath,
		Size:        size,
	}, nil
}

func FileNameToIDAndStarted(filePath string) (*string, *time.Time, error) {
	fileWithExt := filepath.Base(filePath)
	f := strings.TrimSuffix(fileWithExt, filepath.Ext(fileWithExt))
	idAndTime := strings.Split(f, "-")
	if len(idAndTime) != 2 {
		return nil, nil, errors.New("Invalid file name format: " + filePath)
	}
	t, err := time.Parse(TFormatZone, idAndTime[1]+"+0900") // truncated tz +0900
	if err != nil {
		return nil, nil, err
	}
	return &idAndTime[0], &t, nil
}

func (e *Episode) DeleteTSFile(workDir string) error {
	return os.RemoveAll(e.GetTSRawFileDir(workDir))
}

func (e *Episode) DeleteAudioFile(destDir string) error {
	return os.RemoveAll(e.AudioFile(destDir))
}

func (e *Episode) AudioFile(destDir string) string {
	dt := e.Started.Format(TFormat) // truncates tz +0900
	return path.Join(destDir, e.ID+"-"+dt+Encoding)
}

func (e *Episode) ConvertToAudioFile(workDir, destDir string) (string, error) {
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	out := e.AudioFile(destDir)
	err = ffmpeg.Input(e.GetTSRawFilePath(workDir)).
		Output(out, ffmpeg.KwArgs{
			"vn": "",
			"c":  "copy",
			"y":  "",
		}).
		ErrorToStdOut().
		Run()

	if err != nil {
		return "", err
	}
	e.FilePath = out
	return out, nil
}

func (e *Episode) GetTSRawFilePath(workDir string) string {
	return path.Join(e.GetTSRawFileDir(workDir), "video.ts")
}

func (e *Episode) GetTSRawFileDir(workDir string) string {
	return path.Join(workDir, e.ID)
}

func (e *Episode) WriteTag(destDir string) error {
	tag, err := id3v2.Open(e.AudioFile(destDir), id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	tag.SetArtist(e.Program)
	tag.SetAlbum(e.Description)
	tag.SetTitle(e.Title)
	err = tag.Save()
	if err != nil {
		return err
	}

	return nil
}

func (e *Episode) Exists(destDir string) bool {
	_, err := os.Stat(e.AudioFile(destDir))
	return err == nil
}

func (e *Episode) DownloadTS(workDir string) error {
	client := hlsdl.New(e.M3u8Url, nil, e.GetTSRawFileDir(workDir), 64, false)
	_, err := client.Download()
	if err != nil {
		return err
	}
	return nil
}

func (e *Episode) Download(workDir, destDir string) error {
	if e.Exists(destDir) {
		return nil
	}

	defer e.DeleteTSFile(workDir) // nothing to do with the error
	if err := e.DownloadTS(workDir); err != nil {
		return err
	}

	if _, err := e.ConvertToAudioFile(workDir, destDir); err != nil {
		return err
	}

	if err := e.WriteTag(destDir); err != nil {
		if nestedError := e.DeleteAudioFile(destDir); nestedError != nil {
			return nestedError
		}
		return err
	}
	return nil
}
