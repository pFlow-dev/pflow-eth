package service

import (
	"crypto/md5"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"net/http"
	"sync"
	"time"
)

type SessionData struct {
	common.Address
	SessionID string
	LoginAt   time.Time
	LastPing  time.Time
}

type Sessions struct {
	sessions     map[string]*SessionData
	sessionsLock sync.RWMutex
}

func (s *Service) GetSession(sessionID string) (*SessionData, bool) {
	s.sessionsLock.RLock()
	defer s.sessionsLock.RUnlock()
	data, exists := s.sessions[sessionID]
	return data, exists
}

func (s *Service) PingSession(sessionID string, startSession bool) {
	s.sessionsLock.Lock()
	defer s.sessionsLock.Unlock()
	sessionData, exists := s.sessions[sessionID]
	if exists {
		sessionData.LastPing = time.Now()
	} else {
		if startSession {
			s.sessions[sessionID] = &SessionData{
				SessionID: sessionID,
				LastPing:  time.Now(),
				LoginAt:   time.Now(),
			}
		}
	}
}

func (s *Service) CleanupSessions() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.sessionsLock.Lock()
			onlineUsersCount := 0 // Initialize the counter for online users
			for sessionID, sessionData := range s.sessions {
				if now.Sub(sessionData.LastPing) > 1*time.Minute {
					s.Event("session_expire", map[string]interface{}{
						"session":  sessionID,
						"duration": sessionData.LastPing.Sub(sessionData.LoginAt).String(),
						"address":  sessionData.Address,
					})
					delete(s.sessions, sessionID)
				} else {
					onlineUsersCount++ // Increment the counter for online users
				}
			}
			if onlineUsersCount != 0 {
				s.Event("online_users", map[string]interface{}{
					"count": onlineUsersCount,
				})
			}
			s.sessionsLock.Unlock()
		}
	}
}

func (s *Service) TimeOutByAddress(address common.Address) {
	s.sessionsLock.Lock()
	defer s.sessionsLock.Unlock()
	for sessionID, sessionData := range s.sessions {
		if sessionData.Address == address {
			s.Event("session_expire", map[string]interface{}{
				"session": sessionID,
			})
			delete(s.sessions, sessionID)
		}
	}
}

func (s *Service) StartSession(r *http.Request, ts time.Time, address ...common.Address) string {
	session := r.URL.Query().Get("session")
	if session != "" {
		s.PingSession(session, true)
		return session
	} else {
		session = fmt.Sprintf("%x", md5.Sum([]byte(r.RemoteAddr+config.SessionSalt+time.Now().String())))
	}
	s.PingSession(session, true)
	if len(address) > 0 {
		// REVIEW: do we need mutex here?
		s.TimeOutByAddress(address[0]) // Time out the previous session for this address
		s.sessions[session].Address = address[0]
		s.sessions[session].LoginAt = ts
	}
	return session
}

func (s *Service) SetSessionAddress(sessionID string, address common.Address) {
	s.sessionsLock.Lock()
	defer s.sessionsLock.Unlock()
	sessionData, exists := s.sessions[sessionID]
	if exists {
		sessionData.Address = address
	}
}
