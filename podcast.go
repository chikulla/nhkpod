package nhkpod

import (
	"github.com/eduncan911/podcast"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"
)

var Pubdate = time.Date(2011, 9, 1, 0, 0, 0, 0, time.UTC)

func WritePodcastFile(s PodcastSettings, program Program) error {
	p, err := GeneratePodcast(s, program)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.RootDir+"/feed.rss", p.Bytes(), 444)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePodcast(s PodcastSettings, program Program) (*podcast.Podcast, error) {
	audioFiles, err := GetAudioFiles(s.RootDir, program)
	if err != nil {
		return nil, err
	}
	eps := make([]Episode, len(audioFiles))
	for i, f := range audioFiles {
		e, err := NewFromFile(path.Join(s.RootDir, f.Name()), f.Size())
		if err != nil || e.Title == "" {
			continue
		}
		eps[i] = *e
	}
	n := time.Now()
	p := podcast.New(eps[0].Program, "", "", &Pubdate, &n)
	p.Language = "ja_jp"
	p.AddImage(program.ImageUrl)
	for _, ep := range eps {
		p.AddItem(toPodcastItem(s.BaseURL, ep))
	}
	return &p, nil
}

func toPodcastItem(baseUrl string, ep Episode) podcast.Item {
	url := baseUrl + ep.FilePath
	r := podcast.Item{
		GUID:        ep.ID,
		Title:       ep.Title,
		Description: ep.Description,
		PubDate:     &ep.Started,
		Link:        url,
		Enclosure: &podcast.Enclosure{
			URL:    url,
			Type:   podcast.M4A,
			Length: ep.Size,
		},
	}
	return r
}

func GetAudioFiles(audioDir string, program Program) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(audioDir)
	if err != nil {
		return nil, err
	}
	var r []fs.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ex := filepath.Ext(file.Name())
		if ex != Encoding {
			continue
		}
		r = append(r, file)
	}
	return r, nil
}
