package nhkpod

import (
	"encoding/json"
	"net/http"
)

func StartServer(e Env) error {
	assets := http.FileServer(http.Dir(e.AssetsDir))
	audio := http.FileServer(http.Dir(e.AudioDir))
	http.Handle("/", assets)
	http.Handle("/"+e.AudioDir+"/", http.StripPrefix("/"+e.AudioDir+"/", audio))
	http.HandleFunc("/programs", HandlePrograms)
	err := http.ListenAndServe(":"+e.Port, nil)
	if err != nil {
		return err
	}
	return nil
}

func HandlePrograms(w http.ResponseWriter, r *http.Request) {
	p, err := GetAvailablePrograms()
	if err != nil {
		HandleError(w, err)
		return
	}
	b, err := json.Marshal(p)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
