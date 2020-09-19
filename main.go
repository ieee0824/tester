package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/ieee0824/tester/structs"
)

func main() {
	s := flag.String("s", "schedule.toml", "schedule file")
	flag.Parse()

	schedule := &structs.Schedule{}
	if _, err := toml.DecodeFile(*s, schedule); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(schedule.Run())
}
