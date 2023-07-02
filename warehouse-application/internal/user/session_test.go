package user

import (
	"testing"
)

var loginTest = testUser.Login

func TestUser_GenerateSessionKey(t *testing.T) {
	var lastCombo, combo string
	for i := 0; i < 100; i++ {
		combo = generateSessionKey()
		if lastCombo == combo {
			t.Errorf("Repeat combination: %s", combo)
		}
		lastCombo = combo
	}
}

func TestUser_GetSessionKey(t *testing.T) {
	sessionKey, err := getSessionKey(DB, loginTest)
	if sessionKey != "" && err == nil {
		t.Error("Testing login is busy!")
	} else {
		if sessionKey != "" {
			t.Error("Dont correct output session key")
		} else if err == nil {
			t.Error("Dont correct output error")
		}
	}
}

func TestUser_NewSession(t *testing.T) {
	sessionKey, err := getSessionKey(DB, loginTest)
	if sessionKey != "" && err == nil {
		t.Error("Testing login is busy!")
	} else {
		newSessionKey := newSession(DB, loginTest, 0)
		sessionKey, err = getSessionKey(DB, loginTest)
		if err != nil {
			t.Error("Dont correct output error")
		}
		if newSessionKey == "" {
			t.Error("Dont correct output new session key")
		}
		if sessionKey == "" {
			t.Error("Dont correct save new session key in BD")
		}
		if newSessionKey != sessionKey {
			t.Error("New session key dont equals session key in BD")
		}
	}
}

func TestUser_CloseSession(t *testing.T) {
	sessionKey, err := getSessionKey(DB, loginTest)
	if sessionKey == "" || err != nil {
		t.Error("Session dont created!")
	} else {
		testUser.CloseSession(DB)
		sessionKey, err = getSessionKey(DB, loginTest)
		if err == nil {
			t.Error("Dont close session...")
		}
	}
}

func TestUser_StartSession(t *testing.T) {
	sessionKey, err := getSessionKey(DB, loginTest)
	if sessionKey != "" && err == nil {
		t.Error("Testing login is busy!")
	} else {
		newSessionKey := testUser.StartSession(DB)
		if newSessionKey == "" {
			t.Error("Dont correct output new session key")
		} else {
			sessionKey, err := getSessionKey(DB, loginTest)
			if err != nil || sessionKey == "" {
				t.Error("Dont correct save new session key in BD")
			} else if newSessionKey != sessionKey {
				t.Error("New session key dont equals session key in BD")
			}
		}
	}
	testUser.CloseSession(DB)
}
