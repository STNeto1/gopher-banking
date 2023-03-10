// Code generated by ent, DO NOT EDIT.

package ent

import (
	"models/ent/deposit"
	"models/ent/schema"
	"models/ent/transference"
	"models/ent/user"
	"time"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	depositFields := schema.Deposit{}.Fields()
	_ = depositFields
	// depositDescCreatedAt is the schema descriptor for created_at field.
	depositDescCreatedAt := depositFields[3].Descriptor()
	// deposit.DefaultCreatedAt holds the default value on creation for the created_at field.
	deposit.DefaultCreatedAt = depositDescCreatedAt.Default.(func() time.Time)
	// depositDescID is the schema descriptor for id field.
	depositDescID := depositFields[0].Descriptor()
	// deposit.DefaultID holds the default value on creation for the id field.
	deposit.DefaultID = depositDescID.Default.(func() uuid.UUID)
	transferenceFields := schema.Transference{}.Fields()
	_ = transferenceFields
	// transferenceDescCreatedAt is the schema descriptor for created_at field.
	transferenceDescCreatedAt := transferenceFields[4].Descriptor()
	// transference.DefaultCreatedAt holds the default value on creation for the created_at field.
	transference.DefaultCreatedAt = transferenceDescCreatedAt.Default.(func() time.Time)
	// transferenceDescID is the schema descriptor for id field.
	transferenceDescID := transferenceFields[0].Descriptor()
	// transference.DefaultID holds the default value on creation for the id field.
	transference.DefaultID = transferenceDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[5].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
