package main

import (
	"flag"
	"log"
	"os"

	"github.com/ieee0824/tester/structs"
	"gopkg.in/yaml.v2"
)

func main() {
	s := flag.String("s", "schedule.yaml", "schedule file")
	debug := flag.Bool("d", false, "debug mode")
	flag.Parse()

	schedule := &structs.Schedule{}
	if *debug {
		schedule.EnableDebugMode()
	}

	f, err := os.Open(*s)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(schedule); err != nil {
		log.Fatalln(err)
	}

	if err := schedule.Run(); err != nil {
		log.Fatalln(err)
	}

	log.Println("all test is pass")

}
