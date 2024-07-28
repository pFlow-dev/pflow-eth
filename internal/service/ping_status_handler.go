package service

import (
	"encoding/json"
	"fmt"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"net/http"
)

// construct a uuid build string from the JS and CSS build numbers
func getBuild() string {
	hexString := config.JsBuild + "000000000001" + "0000" + config.CssBuild

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexString[0:8], hexString[8:12], hexString[12:16], hexString[16:20], hexString[20:32])
}

var buildUuid = getBuild()

func (s *Service) getSyncData() (map[string]interface{}, error) {
	status := make(map[string]interface{})
	status["sequence"] = s.getContractSequence()
	status["block"] = s.getLatestBlockNumber()

	return status, nil
}

func (s *Service) getSessionData(sessionID string) (map[string]interface{}, error) {
	sessionData, exists := s.GetSession(sessionID)
	status := make(map[string]interface{})
	status["status"] = "ok"
	if exists {
		status["session_data"] = map[string]interface{}{
			"session_id": sessionData.SessionID,
			"last_ping":  sessionData.LastPing,
			"login_at":   sessionData.LoginAt,
		}
	}
	status["network"] = config.Network
	status["build"] = buildUuid
	return status, nil
}

func (s *Service) PingHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	s.PingSession(sessionID, false)
	w.Header().Set("Content-Type", "application/json")
	status, _ := s.getSessionData(sessionID)
	_ = json.NewEncoder(w).Encode(status)
}
