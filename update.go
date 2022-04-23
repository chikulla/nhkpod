package nhkpod

import (
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type PodcastSettings struct {
	ID       string
	CornerID string
	RootDir  string
	WorkDir  string
	BaseURL  string
}

func UpdatePodcasts(e Env) error {
	conf, err := GetConf(e.ConfPath)
	if err != nil {
		return err
	}

	// Clean workdir first
	err = os.RemoveAll(e.WorkDir)
	if err != nil {
		return err
	}

	programs, err := GetAvailablePrograms()
	if err != nil {
		return err
	}
	for _, p := range conf.Podcasts {
		rootDir, err := PreparePodcastRootDir(e, p.ID, p.CornerID)
		setting := PodcastSettings{
			ID:       p.ID,
			CornerID: p.CornerID,
			RootDir:  rootDir,
			WorkDir:  path.Join(e.WorkDir, p.ID+p.CornerID),
			BaseURL:  e.BaseURL,
		}
		if err != nil {
			return err
		}
		err = UpdatePodcast(setting, programs)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdatePodcast(s PodcastSettings, programs *[]Program) error {
	targetProgs, err := GetProgram(s.ID, s.CornerID, programs)
	if err != nil {
		return nil
	}

	err = DownloadProgramEpisodes(s, targetProgs)
	if err != nil {
		return err
	}

	return nil
}

func DownloadProgramEpisodes(s PodcastSettings, progs []Program) error {
	for _, p := range progs {
		err := DownloadEpisodes(s, p)
		if err != nil {
			return err
		}
	}
	if len(progs) > 0 {
		err := WritePodcastFile(s, progs[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadEpisodes(s PodcastSettings, p Program) error {
	eps, err := NewEpisodesFrom(p.URL)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, ep := range *eps {
		wg.Add(1)
		go func(s PodcastSettings, ep Episode) {
			defer wg.Done()
			err := DownloadEpisode(s, ep)
			if err != nil {
				log.Println(err)
			}
		}(s, ep)
	}
	wg.Wait()
	return nil
}

func DownloadEpisode(s PodcastSettings, ep Episode) error {
	err := ep.Download(s.WorkDir, s.RootDir)
	if err != nil {
		return err
	}
	return nil
}

func PreparePodcastRootDir(e Env, programId, cornerId string) (string, error) {
	// If cornerId is not specified, provide all the corners on the program as a single podcast
	id := programId
	if cornerId != "" {
		id = id + "_" + cornerId
	}
	p := path.Join(e.AudioDir, strings.ToLower(id))
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return "", err
	}
	return p, nil
}
