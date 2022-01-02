package account
/*
import (
	"net"
	"strings"
	"time"

	"common/data/entities"
	pass "common/password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manager struct {
	backend DBBackend
}

func NewManager(b DBBackend) *Manager {
	return &Manager{
		backend: b,
	}
}

func (s *Manager) CreateAccount(name string, defaultAdminUsername string, quotaType string, quotaResetTimestamp time.Time, quotaUsd float64, freeBandwidthBytes int64, expiryUtc int64) (primitive.ObjectID, *entities.DefaultUser, error) {
	account, err := entities.NewAccount(name, quotaType, quotaResetTimestamp, quotaUsd, freeBandwidthBytes, expiryUtc)
	if err != nil {
		return account.ID, nil, err
	}

	randomPass := pass.GenerateRandomString(8)
	defaultUser, err := entities.NewUser(account.Name, defaultAdminUsername, randomPass)
	if err != nil {
		return account.ID, nil, err
	}

	//TODO: Workaround --> CHANGE when we implement user verification email
	accID, err := s.backend.CreateWithUser(account, defaultUser)
	if err != nil {
		return account.ID, nil, err
	}
	return accID, &entities.DefaultUser{
		Username: defaultUser.Username,
		Password: randomPass,
	}, err
}

func (s *Manager) AddProject(accountId primitive.ObjectID, name string, quotaUsd float64, vertical string, domain string, proxyType string, isRotating bool, rotationSessionLength time.Duration, whitelist []string, location *string) (*entities.Project, error) {
	project, err := entities.NewProject(name, quotaUsd, vertical, domain, proxyType, isRotating, rotationSessionLength, whitelist, location)
	if err != nil {
		return nil, err
	}
	_, err = s.backend.AddProject(nil, accountId, project)
	if err != nil {
		return nil, err
	}
	return project, nil

}

func (s *Manager) AddRate(accountId primitive.ObjectID, project_type_id primitive.ObjectID, rateUsd float64) (primitive.ObjectID, error) {
	a, err := entities.NewRate(project_type_id, rateUsd)
	if err != nil {
		return a.ID, err
	}
	return s.backend.AddRate(accountId, a)
}

func (s *Manager) GetAccount(id primitive.ObjectID) (*entities.Account, error) {
	return s.backend.Get(id)
}

func (s *Manager) GetAccountByName(name string) (*entities.Account, error) {
	return s.backend.GetByName(strings.ToLower(name), entities.AVAILABLE)
}

func (s *Manager) ListAccounts() (*[]*entities.Account, error) {
	return s.backend.List()
}

func (s *Manager) DeleteAccount(id primitive.ObjectID) error {
	u, err := s.GetAccount(id)
	if err != nil {
		return err
	}
	return s.backend.SoftDelete(u)
}

func (s *Manager) DeleteAccountByName(name string) error {
	u, err := s.GetAccountByName(name)
	if err != nil {
		return err
	}

	return s.backend.SoftDelete(u)
}

func (s *Manager) UpdateAccount(e *entities.Account) error {
	err := e.Validate(true)
	if err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.backend.Update(nil, e)
}

func (s *Manager) IncreaseOrDecreaseAccount(e *entities.Account, proj *entities.Project) error {
	return s.backend.IncreaseOrDecrease(e, proj)
}

func (s *Manager) DeleteProject(accountId primitive.ObjectID, projectID primitive.ObjectID) error {
	return s.backend.DeleteProject(accountId, projectID)
}

func (s *Manager) GetProject(accountId primitive.ObjectID, projectID primitive.ObjectID) (*entities.Project, error) {
	return s.backend.GetProject(nil, accountId, projectID)
}

func (s *Manager) GetProjectByName(accountName string, projectName string) (*entities.Project, error) {
	return s.backend.GetProjectByName(nil, strings.ToLower(accountName), strings.ToLower(projectName))
}

func (s *Manager) UpdateProject(accountId primitive.ObjectID, e *entities.Project) error {
	return s.backend.UpdateProject(accountId, e)
}

func (s *Manager) DeleteRate(accountId primitive.ObjectID, rateID primitive.ObjectID) error {
	return s.backend.DeleteRate(accountId, rateID)
}

func (s *Manager) GetRate(accountId primitive.ObjectID, rateID primitive.ObjectID) (*entities.Rate, error) {
	return s.backend.GetRate(accountId, rateID)
}

func (s *Manager) UpdateRate(accountId primitive.ObjectID, e *entities.Rate) error {
	return s.backend.UpdateRate(accountId, e)
}

func (s *Manager) GetAccountsByIp(ip net.IP) ([]entities.Account, error) {
	return s.backend.GetAccountsByIp(ip)
}

func (s *Manager) AddWhitelistIp(accountId primitive.ObjectID, ip string) error {
	return s.backend.AddWhitelistIp(accountId, ip)
}

func (s *Manager) DeleteWhitelistIp(accountId primitive.ObjectID, ip string) error {
	return s.backend.DeleteWhitelistIp(accountId, ip)
}
*/