package v4

import (
	"time"

	"github.com/jackc/pgtype"
)

// See https://github.com/cyberark/conjur/tree/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate

// Resource v4 resource type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628212349_create_resources.rb
type Resource struct {
	ResourceID string    `gorm:"primaryKey;not null"`
	OwnerID    string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	PolicyID   string    `gorm:"not null"`
}

// Role v4 role type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628212347_create_roles.rb
type Role struct {
	RoleID    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	PolicyID  string    `gorm:"not null"`
}

// Annotation v4 annotation type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628212433_create_annotations.rb
type Annotation struct {
	ResourceID string `gorm:"not null"`
	Name       string `gorm:"not null"`
	Value      string `gorm:"not null"`
	PolicyID   string
}

// Credential v4 credential type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628222441_create_credentials.rb
type Credential struct {
	RoleID        string `gorm:"not null"`
	ClientID      string
	APIKey        []byte `gorm:"type:bytea"`
	EncryptedHash []byte `gorm:"type:bytea"`
	Expiration    time.Time
}

// HostFactoryToken v4 host_factory_token type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20170404125612_create_host_factories.rb
type HostFactoryToken struct {
	TokenSha256 string        `gorm:"type:varchar(64);not null;"`
	Token       []byte        `gorm:"type:bytea;not null"`
	ResourceID  string        `gorm:"not null"`
	Cidr        []pgtype.Inet `gorm:"type:cidr[];not null"`
	Expiration  time.Time
}

// Permission v4 permissions type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628212428_create_permissions.rb
type Permission struct {
	Privilege  string `gorm:"not null"`
	ResourceID string `gorm:"not null"`
	RoleID     string `gorm:"not null"`
	PolicyID   string
}

// PolicyVersion v4 policy_versions type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160815131453_create_policy_version.rb
type PolicyVersion struct {
	ResourceID string `gorm:"not null"`
	RoleID     string `gorm:"not null"`
	Version    int    `gorm:"not null;type:integer"`
}

// ResourcesTextsearch v4 resource_textsearch type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20170710163523_create_resources_textsearch.rb
type ResourcesTextsearch struct {
	ResourceID string `gorm:"not null"`
	Textsearch string `gorm:"type:tsvector"`
}

// TableName returns the table name to use for ResourcesTextsearch
func (ResourcesTextsearch) TableName() string {
	return "resources_textsearch"
}

// RoleMembership v4 role_memberships type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160628212358_create_role_memberships.rb
type RoleMembership struct {
	RoleID      string `gorm:"not null"`
	MemberID    string `gorm:"not null"`
	AdminOption bool   `gorm:"not null"`
	Ownership   bool   `gorm:"not null"`
	PolicyID    string
}

// SchemaMigration v4 schema_migrations type
type SchemaMigration struct {
	Filename string `gorm:"not null;primaryKey"`
}

// Secret v4 secrets type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20160630172059_create_secrets.rb
type Secret struct {
	ResourceID string `gorm:"not null"`
	Version    int    `gorm:"not null;type:integer"`
	Value      []byte `gorm:"not null;type:bytea"`
}

// SlosiloKeystore v4 slosilo_keystor type
// https://github.com/cyberark/conjur/blob/f12d4522ee1ef175715fc01c2eb723dc18922f8e/db/migrate/20121215032820_create_keystore.rb
type SlosiloKeystore struct {
	ID          string `gorm:"not null"`
	Key         []byte `gorm:"not null;type:bytea"`
	Fingerprint string `gorm:"not null"`
}

// TableName v4 slosilo_keystore table name
func (SlosiloKeystore) TableName() string {
	return "slosilo_keystore"
}

// Data a full database of data
type Data struct {
	Resources           []Resource
	Roles               []Role
	Annotations         []Annotation
	Credentials         []Credential
	HostFactoryTokens   []HostFactoryToken
	Permissions         []Permission
	PolicyVersions      []PolicyVersion
	ResourcesTextsearch []ResourcesTextsearch
	RoleMemberships     []RoleMembership
	SchemaMigrations    []SchemaMigration
	Secrets             []Secret
	SlosiloKeystore     []SlosiloKeystore
}
