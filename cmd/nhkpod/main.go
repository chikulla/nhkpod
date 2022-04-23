package main

import (
	nhkpod2 "github.com/chikulla/nhkpod"
	"log"
	"os"
)

func main() {
	env, err := setup()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("initial audio download")
	if err := nhkpod2.UpdatePodcasts(*env); err != nil {
		log.Println(err)
	}

	if err := nhkpod2.StartScheduler(*env); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("server starting on: ", env.Port)
	if err := nhkpod2.StartServer(*env); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func setup() (*nhkpod2.Env, error) {
	env := nhkpod2.Env{}
	if err := env.Init(); err != nil {
		return nil, err
	}
	if err := nhkpod2.SetupLogOutput(env); err != nil {
		return nil, err
	}
	return &env, nil
}
