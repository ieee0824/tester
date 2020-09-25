package structs

import "fmt"

type Schedule struct {
	debug bool
	Jobs  []Job
}

func (s *Schedule) EnableDebugMode() {
	s.debug = true
}

func (s *Schedule) DisableDebugMode() {
	s.debug = false
}

func (s *Schedule) Run() error {
	var opt []JobOption
	if s.debug {
		opt = []JobOption{
			{Debug: true},
		}
	}
	for _, j := range s.Jobs {
		if err := j.Run(opt...); err != nil {
			return fmt.Errorf("Error: job name: %s, error: %v", j.Name, err)
		}
	}

	return nil
}
