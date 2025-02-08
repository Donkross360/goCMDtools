package cmd

import (
	"interactiveTools/pomo/pomodoro"
	"interactiveTools/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
