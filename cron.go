package nhkpod

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

func StartScheduler(e Env) error {
	tz, err := time.LoadLocation("JST")
	if err != nil {
		tz = time.UTC
	}
	s := gocron.NewScheduler(tz)

	_, err = s.Cron(e.Schedule).SingletonMode().Do(func() {
		err := UpdatePodcasts(e)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}
