package nhkpod

import (
	"io"
	"log"
	"os"
)

func SetupLogOutput(e Env) error {
	file, err := os.OpenFile(e.LogFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	w := io.MultiWriter(os.Stdout, file)
	log.SetOutput(w)
	return nil
}
