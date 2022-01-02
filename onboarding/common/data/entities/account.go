package entities
/*
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

var AccountUniqueKeys = []string{"name", "_id"}

type Status int

const (
	AVAILABLE Status = iota + 1
	DISABLED
	SUSPENDED
	DELETED
)

type Account struct {
	Entity                 `bson:"-"`
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt              time.Time          `bson:"created_at,omitempty"`
	UpdatedAt              time.Time          `bson:"updated_at,omitempty"`
	Name                   string             `bson:"name,omitempty"`
	Status                 Status             `bson:"status,omitempty"`
	QuotaType              string             `bson:"quota_type,omitempty"`
	QuotaResetTimestamp    time.Time          `bson:"quota_reset_timestamp,omitempty"`
	FreeBandwidthBytes     *int64             `bson:"free_bandwidth_bytes,omitempty"`
	FreeBandwidthUsedBytes *int64             `bson:"free_bandwidth_used_bytes,omitempty"`
	ExpiryUtc              *int64             `bson:"expiry_utc,omitempty"`
	QuotaUsd               *float64           `bson:"quota_usd,omitempty"`
	SpentUsd               *float64           `bson:"spent_usd,omitempty"`
	SpentBytes             *int64             `bson:"spent_bytes,omitempty"`
	GUID                   string             `bson:"guid,omitempty"`
	QuotaGraceDays         *float64           `bson:"quota_grace_days,omitempty"`
	Projects               []Project          `bson:"projects,omitempty" flatupdate:"-"`
	Rates                  []Rate             `bson:"rates,omitempty" flatupdate:"-"`
	WhitelistIps           []string           `bson:"whitelist_ips" flatupdate:"-"`
}

func NewAccount(name string, quotaType string, quotaResetTimestamp time.Time, quotaUsd float64, freeBandwidthBytes int64, expiryUtc int64) (*Account, error) {
	now := time.Now()
	int64Zero := int64(0)
	float64Zero := float64(0)

	account := &Account{
		ID:                     primitive.NewObjectID(),
		CreatedAt:              now,
		UpdatedAt:              now,
		Name:                   strings.ToLower(name),
		QuotaType:              quotaType,
		QuotaResetTimestamp:    quotaResetTimestamp,
		QuotaUsd:               &quotaUsd,
		FreeBandwidthBytes:     &freeBandwidthBytes,
		FreeBandwidthUsedBytes: &int64Zero,
		ExpiryUtc:              &expiryUtc,
		SpentUsd:               &float64Zero,
		SpentBytes:             &int64Zero,
		QuotaGraceDays:         &float64Zero,
		Projects:               []Project{},
		Rates:                  []Rate{},
		WhitelistIps:           []string{},
		Status:                 AVAILABLE,
	}

	err := account.Validate(false)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *Account) Validate(partial bool) error {
	return nil
}
*/