package nhkpod

import (
	"encoding/json"
	"errors"
	"github.com/grafov/m3u8"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ProgramDetailJson struct {
	Main struct {
		SiteID       string `json:"site_id"`
		ProgramTitle string `json:"program_name"`
		CornerID     string `json:"corner_id"`
		CornerTitle  string `json:"corner_name"`
		DetailList   []struct {
			FileList []struct {
				FileID       string `json:"file_id"`
				FileTitle    string `json:"file_title"`
				FileTitleSub string `json:"file_title_sub"`
				FileName     string `json:"file_name"`
				AAVInfo1     string `json:"aa_vinfo1"`
				AAVInfo4     string `json:"aa_vinfo4"`
			} `json:"file_list"`
		} `json:"detail_list"`
	} `json:"main"`
}

func NewEpisodesFrom(programURL string) (*[]Episode, error) {
	r, err := http.Get(programURL)
	if err != nil {
		return nil, err
	}

	var programDetail = &ProgramDetailJson{}
	b, err := ioutil.ReadAll(r.Body)
	if err = json.Unmarshal(b, programDetail); err != nil {
		return nil, err
	}

	numOfEp := 0
	for _, d := range programDetail.Main.DetailList {
		numOfEp += len(d.FileList)
	}
	i := 0
	ep := make([]Episode, numOfEp)
	for _, d := range programDetail.Main.DetailList {
		for _, f := range d.FileList {
			started, err := ConvertFileTime(f.AAVInfo4)
			if err != nil {
				continue
			}
			url, err := GetM3U8URL(f.FileName)
			if err != nil {
				continue
			}
			desc := f.AAVInfo1 + " " + f.FileTitleSub
			ep[i] = Episode{
				ID:          f.FileID,
				M3u8Url:     *url,
				Title:       f.FileTitle,
				Description: desc,
				Program:     strings.Join([]string{programDetail.Main.ProgramTitle, programDetail.Main.CornerTitle}, " "),
				Started:     *started,
			}
			i++
		}
	}
	return &ep, err
}

func GetM3U8URL(url string) (*string, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	playlist, listType, err := m3u8.DecodeFrom(r.Body, true)
	if err != nil {
		return nil, err
	}
	if listType != m3u8.MASTER {
		return nil, errors.New("m3u8 list type invalid")
	}
	p := playlist.(*m3u8.MasterPlaylist)

	if p == nil || len(p.Variants) < 1 {
		return nil, errors.New("m3u8 format invalid")
	}
	return &p.Variants[0].URI, nil
}

func ConvertFileTime(info04 string) (*time.Time, error) {
	// parsing text like 2022-04-22T00:50:00+09:00_2022-04-22T01:00:00+09:00
	start := strings.Split(info04, "_")[0]
	parsed, err := time.Parse("2006-01-02T15:04:05-07:00", start)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

type TimeTableJson struct {
	DataList []struct {
		SiteID      string `json:"site_id"`
		ProgramName string `json:"program_name"`
		DetailJson  string `json:"detail_json"`
		CornerID    string `json:"corner_id"`
		CornerName  string `json:"corner_name"`
		ThumbnailP  string `json:"thumbnail_p"`
	} `json:"data_list"`
}

func GetTimeTable() (*TimeTableJson, error) {
	var result = &TimeTableJson{}
	url := "https://www.nhk.or.jp/radioondemand/json/index_v3/index.json"
	r, err := http.Get(url)
	if err != nil {
		return result, err
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return result, err
	}
	return result, err
}

type Program struct {
	ID       string
	CornerID string
	ImageUrl string
	URL      string
	Title    string
	Corner   string
}

func GetAvailablePrograms() (*[]Program, error) {
	j, err := GetTimeTable()
	if err != nil {
		return nil, err
	}
	pgs := make([]Program, len(j.DataList))
	for i, d := range j.DataList {
		pgs[i] = Program{
			URL:      d.DetailJson,
			ID:       d.SiteID,
			ImageUrl: d.ThumbnailP,
			Title:    d.ProgramName,
			CornerID: d.CornerID,
			Corner:   d.CornerName,
		}
	}
	return &pgs, err
}

func GetProgram(id, cornerID string, programs *[]Program) ([]Program, error) {
	var err error
	if programs == nil {
		programs, err = GetAvailablePrograms()
		if err != nil {
			return nil, err
		}
	}
	var result []Program
	for _, p := range *programs {
		if id == p.ID {
			if cornerID == "" || cornerID != "" && cornerID == p.CornerID {
				result = append(result, p)
			}
		}
	}

	return result, nil
}
