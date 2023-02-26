// Code generated by ent, DO NOT EDIT.

package transference

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the transference type in the database.
	Label = "transference"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAmount holds the string denoting the amount field in the database.
	FieldAmount = "amount"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeFromUser holds the string denoting the from_user edge name in mutations.
	EdgeFromUser = "from_user"
	// EdgeToUser holds the string denoting the to_user edge name in mutations.
	EdgeToUser = "to_user"
	// Table holds the table name of the transference in the database.
	Table = "transferences"
	// FromUserTable is the table that holds the from_user relation/edge.
	FromUserTable = "transferences"
	// FromUserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	FromUserInverseTable = "users"
	// FromUserColumn is the table column denoting the from_user relation/edge.
	FromUserColumn = "user_from_transfers"
	// ToUserTable is the table that holds the to_user relation/edge.
	ToUserTable = "transferences"
	// ToUserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	ToUserInverseTable = "users"
	// ToUserColumn is the table column denoting the to_user relation/edge.
	ToUserColumn = "user_to_transfers"
)

// Columns holds all SQL columns for transference fields.
var Columns = []string{
	FieldID,
	FieldAmount,
	FieldMessage,
	FieldStatus,
	FieldCreatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "transferences"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_from_transfers",
	"user_to_transfers",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// Status defines the type for the "status" enum field.
type Status string

// Status values.
const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusDenied    Status = "denied"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPending, StatusCompleted, StatusDenied:
		return nil
	default:
		return fmt.Errorf("transference: invalid enum value for status field: %q", s)
	}
}