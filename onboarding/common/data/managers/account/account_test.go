package account
/*
import (
	"errors"
	"testing"
	"time"

	"common/data/dbbackends"
	"common/data/dbbackends/mongo"
	"common/data/dbbackends/mongo/account"
	"common/data/entities"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	AccountManager := NewManager(account.NewMongoAccount(mongo.NewMongoConnector()))
	id, user, err := AccountManager.CreateAccount("PROWLERCORP", "username@valid.com", "none", time.Now(), 2, 0, 0)
	if err != nil {
		t.Fatalf("failed creating Account with err: %v", err)
	}
	Account, err := AccountManager.GetAccount(id)
	assert.Equal(t, id, Account.ID)
	assert.Equal(t, "prowlercorp", Account.Name)
	assert.Equal(t, "none", Account.QuotaType)
	assert.NotNil(t, user)

	t.Log(Account)
	if err != nil {
		t.Fatalf("failed getting Account with err: %v", err)
	}
	err = AccountManager.UpdateAccount(&entities.Account{ID: Account.ID, QuotaType: "fixed"})
	if err != nil {
		t.Fatalf("failed updating creating Account with err: %v", err)
	}
	Accounts, err := AccountManager.ListAccounts()
	assert.Equal(t, 1, len(*Accounts))
	for _, p := range *Accounts {
		t.Log(p)
		assert.Equal(t, "fixed", p.QuotaType)
		assert.Equal(t, "prowlercorp", p.Name)

	}
	if err != nil {
		t.Fatalf("failed listing Accounts with err: %v", err)
	}

	spentBytes := int64(2)
	spentUsd := float64(-2)
	err = AccountManager.IncreaseOrDecreaseAccount(&entities.Account{ID: Account.ID, SpentBytes: &spentBytes, SpentUsd: &spentUsd}, nil)
	if err != nil {
		t.Fatalf("failed increasing quota of account with err: %v", err)
	}
	Account, err = AccountManager.GetAccount(id)
	t.Log(Account)
	assert.Equal(t, int64(2), *Account.SpentBytes)
	assert.Equal(t, float64(-2), *Account.SpentUsd)

	err = AccountManager.IncreaseOrDecreaseAccount(&entities.Account{ID: Account.ID, QuotaType: "aa"}, nil)
	if err == nil {
		t.Fatalf("expected increase function to fail with bad input")
	}

	err = AccountManager.DeleteAccount(id)
	if err != nil {
		t.Fatalf("failed deleting Account with err: %v", err)
	}
	Account, err = AccountManager.GetAccount(id)
	t.Log(Account)
	if errors.Is(err, &dbbackends.ErrNotFound{}) {
		t.Fatal("failed getting Account after deletion, should have been not found error")
	}
}

func TestProject(t *testing.T) {
	AccountManager := NewManager(account.NewMongoAccount(mongo.NewMongoConnector()))
	id, user, err := AccountManager.CreateAccount("PROWLERCORP", "username@valid.com", "none", time.Now(), 2, 0, 0)
	if err != nil {
		t.Fatalf("failed creating Account with err: %v", err)
	}
	Account, err := AccountManager.GetAccount(id)
	assert.Equal(t, 0, len(Account.Projects))
	assert.Equal(t, int64(0), *Account.ExpiryUtc)
	assert.NotNil(t, user)

	t.Log(Account)
	if err != nil {
		t.Fatalf("failed getting Account with err: %v", err)
	}

	invalidPojectNames := []string{
		"as",
		"invalidVeryLongProjectNameForBadRequest",
		"invalid_name",
		"invalidName!",
	}
	noCountry := ""
	for _, name := range invalidPojectNames {
		t.Log("creating invalid project name: ", name)
		_, err := AccountManager.AddProject(Account.ID, name, 5, "serp", "google", "residential", true, 0, nil, &noCountry)
		assert.Equal(t, entities.ErrInvalidEntity, err)
	}

	validName := "test1"
	p, err := AccountManager.AddProject(Account.ID, validName, 5, "serp", "google", "residential", true, 0, nil, &noCountry)
	if err != nil {
		t.Fatalf("failed creating project with err: %v", err)
	}
	Account, err = AccountManager.GetAccount(id)
	assert.Equal(t, validName, Account.Projects[0].Name)

	err = AccountManager.UpdateProject(Account.ID, &entities.Project{ID: p.ID, QuotaUsd: 10})
	if err != nil {
		t.Fatalf("failed updating project with err: %v", err)
	}
	Account, err = AccountManager.GetAccount(id)
	t.Log(Account)
	assert.Equal(t, float64(10), Account.Projects[0].QuotaUsd)
	assert.Equal(t, validName, Account.Projects[0].Name)

	Project, err := AccountManager.GetProject(id, p.ID)
	t.Log(Project)
	assert.Equal(t, float64(10), Project.QuotaUsd)

	Project, err = AccountManager.GetProjectByName(Account.Name, validName)
	t.Log(Project)
	assert.Equal(t, validName, Project.Name)

	err = AccountManager.DeleteProject(Account.ID, p.ID)
	if err != nil {
		t.Fatalf("failed deleting project with err: %v", err)
	}
	Account, err = AccountManager.GetAccount(id)
	t.Log(Account)
	assert.Equal(t, 0, len(Account.Projects))
}
*/