package oauth2_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/sirupsen/logrus"
)

var clientManagers = map[string]pkg.FositeStorer{}
var clientManager = &client.MemoryManager{
	Clients: map[string]client.Client{"foobar": {ID: "foobar"}},
	Hasher:  &fosite.BCrypt{},
}

func init() {
	clientManagers["memory"] = &FositeMemoryStore{
		AuthorizeCodes: make(map[string]fosite.Requester),
		IDSessions:     make(map[string]fosite.Requester),
		AccessTokens:   make(map[string]fosite.Requester),
		RefreshTokens:  make(map[string]fosite.Requester),
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		connectToPG()
		connectToMySQL()
		connectToPGConsent()
		connectToMySQLConsent()
	}

	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New()}
	if _, err := s.CreateSchemas(); err != nil {
		logrus.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New()}
	if _, err := s.CreateSchemas(); err != nil {
		logrus.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func TestCreateGetDeleteAuthorizeCodes(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAuthorizeCodes(m))
	}
}

func TestCreateGetDeleteAccessTokenSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAccessTokenSession(m))
	}
}

func TestCreateGetDeleteOpenIDConnectSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteOpenIDConnectSession(m))
	}
}

func TestCreateGetDeleteRefreshTokenSession(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteRefreshTokenSession(m))
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	t.Parallel()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperRevokeRefreshToken(m))
	}
}
