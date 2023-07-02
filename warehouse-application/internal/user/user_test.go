package user

import (
	"testing"
	config2 "warehouse-application/internal/config"
	"warehouse-application/internal/postgresSQL"
)

var testUser = User{"TestLogin123456789", "TestPassword123456789"}
var configPath = "/Users/kare/GolandProjects/warehouse-application/config/config.yaml"
var config = config2.MustLoad(configPath)
var dataSourceName = config.PostgresConfig.GetDataSourceName()
var DB = postgresSQL.Connect(dataSourceName)

func TestDontHaveUser(t *testing.T) {
	hashPassword, err := testUser.getHashPassword(DB)
	if hashPassword != "" {
		t.Errorf("User have hash password: %s", hashPassword)
	}
	if err == nil {
		t.Error("Dont have error")
	}
}

func TestRegistrationUser(t *testing.T) {
	err := testUser.RegisterUser(DB)
	if err != nil {
		t.Error("User already registered!")
	} else {
		hashPassword, err := testUser.getHashPassword(DB)
		if hashPassword == "" {
			t.Error("User dont have hash password")
		}
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}
}

func TestLoginUser(t *testing.T) {
	keySession, err := testUser.LoginUser(DB)
	if err != nil {
		t.Errorf("Dont correct output error: %v", err)
	}
	if keySession == "" {
		t.Error("Empty key session")
	}
}

func TestSuccessfulExitUser(t *testing.T) {
	err := testUser.ExitUser(DB)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestErrorExitUser(t *testing.T) {
	err := testUser.ExitUser(DB)
	if err == nil {
		t.Error("Dont get error")
	}
}

func TestRemoveUser(t *testing.T) {
	err := testUser.removeUser(DB)
	if err != nil {
		t.Errorf("Erorr remove user: %v", err)
	}
}

func TestLoginEmptyUser(t *testing.T) {
	keySession, err := testUser.LoginUser(DB)
	if keySession != "" || err == nil {
		t.Error("Dont correct output error empty user")
	}
}
