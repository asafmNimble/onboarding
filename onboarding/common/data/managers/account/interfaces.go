package account
/*
import (
	"errors"
	"net"
	"time"

	"common/data/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Reader interface {
	Get(id primitive.ObjectID) (*entities.Account, error)
	GetByName(name string, status entities.Status) (*entities.Account, error)
	List() (*[]*entities.Account, error)
	GetRate(accountId primitive.ObjectID, rateID primitive.ObjectID) (*entities.Rate, error)
	GetProject(session mongo.SessionContext, accountId primitive.ObjectID, projectID primitive.ObjectID) (*entities.Project, error)
	GetProjectByName(session mongo.SessionContext, accountName string, projectName string) (*entities.Project, error)
}

type Writer interface {
	//CreateWithProject(e *entities.Account) (primitive.ObjectID, error)
	Create(session mongo.SessionContext, e *entities.Account) (primitive.ObjectID, error)
	CreateWithUser(a *entities.Account, u *entities.User) (primitive.ObjectID, error)
	Update(session mongo.SessionContext, e *entities.Account) error
	Delete(id primitive.ObjectID) error
	DeleteByName(name string) error
	SoftDelete(e *entities.Account) error
	IncreaseOrDecrease(e *entities.Account, proj *entities.Project) error

	AddProject(session mongo.SessionContext, accountId primitive.ObjectID, e *entities.Project) (primitive.ObjectID, error)
	DeleteProject(accountId primitive.ObjectID, projectID primitive.ObjectID) error
	UpdateProject(accountId primitive.ObjectID, e *entities.Project) error

	AddRate(accountId primitive.ObjectID, e *entities.Rate) (primitive.ObjectID, error)
	DeleteRate(accountId primitive.ObjectID, rateID primitive.ObjectID) error
	UpdateRate(accountId primitive.ObjectID, e *entities.Rate) error

	GetAccountsByIp(ip net.IP) ([]entities.Account, error)
	AddWhitelistIp(accountId primitive.ObjectID, ip string) error
	DeleteWhitelistIp(accountId primitive.ObjectID, ip string) error
}

type DBBackend interface {
	Reader
	Writer
}

type AccountManager interface {
	GetAccount(id primitive.ObjectID) (*entities.Account, error)
	GetAccountByName(accountName string) (*entities.Account, error)
	ListAccounts() (*[]*entities.Account, error)
	CreateAccount(name string, defaultAdminUsername string, quotaType string, quotaResetTimestamp time.Time, quotaUsd float64, freeBandwidthBytes int64, expiryUtc int64) (primitive.ObjectID, *entities.DefaultUser, error)
	UpdateAccount(e *entities.Account) error
	DeleteAccount(id primitive.ObjectID) error
	DeleteAccountByName(name string) error
	IncreaseOrDecreaseAccount(e *entities.Account, proj *entities.Project) error

	AddProject(accountId primitive.ObjectID, name string, quotaUsd float64, vertical string, domain string, proxyType string, isRotating bool, RotationSessionlength time.Duration, whitelist []string, location string) (*entities.Project, error)
	DeleteProject(accountId primitive.ObjectID, projectID primitive.ObjectID) error
	UpdateProject(accountId primitive.ObjectID, e *entities.Project) error
	GetProject(accountId primitive.ObjectID, projectID primitive.ObjectID) (*entities.Project, error)
	GetProjectByName(accountName string, projectName string) (*entities.Project, error)

	AddRate(accountId primitive.ObjectID, project_type_id primitive.ObjectID, rateUsd float64) (primitive.ObjectID, error)
	DeleteRate(accountId primitive.ObjectID, rateID primitive.ObjectID) error
	UpdateRate(accountId primitive.ObjectID, e *entities.Rate) error
	GetRate(accountId primitive.ObjectID, rateID primitive.ObjectID) (*entities.Rate, error)

	GetAccountsByIp(ip net.IP) ([]entities.Account, error)
	AddWhitelistIp(accountId primitive.ObjectID, ip string) error
	DeleteWhitelistIp(accountId primitive.ObjectID, ip string) error
}

var ErrDuplicateAccount = errors.New("Account with that name already exists")
var ErrDuplicateUser = errors.New("User with that name already exists")
*/