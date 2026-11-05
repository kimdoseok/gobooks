package sv

import (
	"sync"
	"time"
)

type (
	AppState struct {
		mu sync.RWMutex

		theme       string
		loginStatus bool
		loginUser   User
	}
	User struct {
		ID         int64
		Yus_code   string
		Yus_passwd string
		Yus_name   string
		Yus_memo   string
		Yus_active bool
		Yus_system bool
		CreatedAt  *time.Time
		UpdatedAt  *time.Time
	}
)

var Global = &AppState{
	theme: "light",
}

func (s *AppState) Theme() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.theme
}

func (s *AppState) SetTheme(theme string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.theme = theme
}

func (s *AppState) LoginStatus() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loginStatus
}

func (s *AppState) SetLoginStatus(status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.loginStatus = status
}

func (s *AppState) LoginUser() User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loginUser
}

func (s *AppState) SetLoginUser(user User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.loginUser = user
}
