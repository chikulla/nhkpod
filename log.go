package nhkpod

import (
	"io"
	"log"
	"os"
)

func SetupLogOutput(e Env) error {
	err := os.Remove(e.LogFile)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(e.LogFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	w := io.MultiWriter(os.Stdout, file)
	log.SetOutput(w)
	return nil
}
