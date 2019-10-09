package integration

import (
	"testing"

	"chainup.dev/chainup"
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/httpserver/routes"
	"chainup.dev/lib/test"
	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"
)

// TestConfig used for running integration tests.
type TestConfig struct {
	TargetAPI string `yaml:"target_api"`
}

// Validate satisfies the config.Config interface.
func (cfg *TestConfig) Validate() error {
	if cfg.TargetAPI == "" {
		return errors.New("missing target API url")
	}
	return nil
}

func TestAuthenticationFlow(t *testing.T) {
	initEnvironment(t)

	email := account.NewEmail(randomdata.Email())
	password := account.NewClearPassword(randomdata.SillyName())

	// Check if user can register an account.
	acc, token := registerAccount(t, email, password)
	test.AssertStringsEqual(t, "email matches", acc.Email.String(), email.String())
	test.AssertStringsEqual(t, "password not returned", acc.Password.String(), "")

	// Check if user can access a protected endpoint with register token.
	authenticateAs(token)
	getAccount(t, "me")
	getAccount(t, acc.ID.String())

	// Check if user can login into an account.
	token = loginAccount(t, email, password)

	// Check if user can access a protected endpoint with login token.
	authenticateAs(token)
	getAccount(t, "me")
	getAccount(t, acc.ID.String())
}

func TestBadRegistrationFlow(t *testing.T) {
	initEnvironment(t)

	err := test.SendPost("/register", &routes.RegisterRequest{}, 400, nil)
	test.CheckErr(t, "send invalid register request", err)
}

func TestBadLoginFlow(t *testing.T) {
	initEnvironment(t)

	err := test.SendPost("/login", &routes.LoginRequest{}, 400, nil)
	test.CheckErr(t, "send invalid login request", err)

	err = test.SendPost("/login", &routes.LoginRequest{
		Email:    "invalid@example.com",
		Password: "wrongpass",
	}, 403, nil)
	test.CheckErr(t, "send invalid login request with email", err)
}

func initEnvironment(t *testing.T) {
	test.Integration(t)

	var cfg TestConfig
	_, err := chainup.ProvideTestConfigProvider().Load(&cfg)
	test.CheckErr(t, "initialize config", err)

	test.SetBaseURL(cfg.TargetAPI)
}

func registerAccount(t *testing.T, email account.Email, password account.ClearPassword) (*account.Account, account.Token) {
	var registerReq = routes.RegisterRequest{Email: email, Password: password}
	var registerResp routes.RegisterResponse

	err := test.SendPost("/register", &registerReq, 201, &registerResp)
	test.CheckErr(t, "send register request", err)

	return registerResp.Account, registerResp.Token
}

func loginAccount(t *testing.T, email account.Email, password account.ClearPassword) account.Token {
	var loginReq = routes.LoginRequest{Email: email, Password: password}
	var loginResp routes.LoginResponse

	err := test.SendPost("/login", &loginReq, 200, &loginResp)
	test.CheckErr(t, "send login request", err)

	return loginResp.Token
}

func authenticateAs(token account.Token) {
	test.SetHeader("Authorization", "Bearer "+token.String())
}

func getAccount(t *testing.T, id string) *account.Account {
	var acc account.Account

	err := test.SendGet("api/v1/account/"+id, 200, &acc)
	test.CheckErr(t, "get account by ID", err)

	return &acc
}
