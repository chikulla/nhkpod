package nhkpod

import (
	"net"
	"os"
)

type Env struct {
	Schedule  string //cron format
	LogFile   string
	AssetsDir string
	AudioDir  string
	WorkDir   string
	ConfPath  string
	Port      string
	Host      string
	Https     bool
	BaseURL   string
}

func (e *Env) InitBaseURl() {
	if e.BaseURL != "" {
		return
	}
	protocol := "http://"
	if e.Https {
		protocol = "https://"
	}
	e.BaseURL = protocol + e.Host + ":" + e.Port + "/"
}

func (e *Env) ReadEnv() {
	e.Schedule = GetenvOrDefault(e.Schedule, "NHKPOD_SCHEDULE", "35 * * * *")
	e.LogFile = GetenvOrDefault(e.LogFile, "NHKPOD_LOG_FILE", "log.log")
	e.AssetsDir = GetenvOrDefault(e.AssetsDir, "NHKPOD_ASSETS_DIR", "assets")
	e.AudioDir = GetenvOrDefault(e.AudioDir, "NHKPOD_AUDIO_DIR", "audio")
	e.WorkDir = GetenvOrDefault(e.WorkDir, "NHKPOD_WORK_DIR", "temp")
	e.ConfPath = GetenvOrDefault(e.ConfPath, "NHKPOD_CONF_PATH", "conf.yml") // in case need to get conf via http
	e.Port = GetenvOrDefault(e.Port, "NHKPOD_PORT", "8080")
	e.Host = GetenvOrDefault(e.Host, "NHKPOD_HOST", "") // InitHost would resolve

	if e.Https {
		return
	}

	if os.Getenv("NHKPOD_HTTPS") != "" {
		e.Https = true
	}
}

func GetenvOrDefault(exists, key, defaultVal string) string {
	if exists != "" {
		return exists
	}
	e := os.Getenv(key)
	if e == "" {
		e = defaultVal
	}
	return e
}

func (e *Env) InitHost() error {
	if e.Host == "" {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return err
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		e.Host = localAddr.IP.String()
	}
	return nil
}

func (e *Env) Init() error {
	e.ReadEnv()
	if err := e.InitHost(); err != nil {
		return err
	}
	e.InitBaseURl()
	if err := SetupDir([]string{e.WorkDir, e.AssetsDir}); err != nil {
		return err
	}
	return nil
}

func SetupDir(dirs []string) error {
	for _, d := range dirs {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
