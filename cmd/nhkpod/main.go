package main

import (
	nhkpod2 "github.com/chikulla/nhkpod"
	"log"
	"os"
)

func main() {
	env, conf, err := setup()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("initial audio download")
	if err := nhkpod2.UpdatePodcasts(*env, *conf); err != nil {
		log.Println(err)
	}

	if err := nhkpod2.StartScheduler(*env, *conf); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("server starting on: ", env.Port)
	if err := nhkpod2.StartServer(*env); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func setup() (*nhkpod2.Env, *nhkpod2.Conf, error) {
	env := nhkpod2.Env{}
	if err := env.Init(); err != nil {
		return nil, nil, err
	}
	if err := nhkpod2.SetupLogOutput(env); err != nil {
		return nil, nil, err
	}

	conf, err := nhkpod2.GetConf(env.ConfPath)
	if err != nil {
		return nil, nil, err
	}

	return &env, conf, nil
}
