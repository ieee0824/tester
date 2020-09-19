package structs

import "fmt"

type Schedule struct {
	Jobs []Job
}

func (s *Schedule) Run() error {
	for _, j := range s.Jobs {
		if err := j.Run(); err != nil {
			return fmt.Errorf("Error: job name: %s, error: %v", j.Name, err)
		}
	}

	return nil
}
