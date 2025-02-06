package pomodoro_test

import (
	"testing"

	"interactiveTools/pomo/pomodoro"
	"interactiveTools/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
