// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"models/ent/deposit"
	"models/ent/predicate"
	"models/ent/user"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DepositUpdate is the builder for updating Deposit entities.
type DepositUpdate struct {
	config
	hooks    []Hook
	mutation *DepositMutation
}

// Where appends a list predicates to the DepositUpdate builder.
func (du *DepositUpdate) Where(ps ...predicate.Deposit) *DepositUpdate {
	du.mutation.Where(ps...)
	return du
}

// SetAmount sets the "amount" field.
func (du *DepositUpdate) SetAmount(f float64) *DepositUpdate {
	du.mutation.ResetAmount()
	du.mutation.SetAmount(f)
	return du
}

// AddAmount adds f to the "amount" field.
func (du *DepositUpdate) AddAmount(f float64) *DepositUpdate {
	du.mutation.AddAmount(f)
	return du
}

// SetStatus sets the "status" field.
func (du *DepositUpdate) SetStatus(d deposit.Status) *DepositUpdate {
	du.mutation.SetStatus(d)
	return du
}

// SetCreatedAt sets the "created_at" field.
func (du *DepositUpdate) SetCreatedAt(t time.Time) *DepositUpdate {
	du.mutation.SetCreatedAt(t)
	return du
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (du *DepositUpdate) SetNillableCreatedAt(t *time.Time) *DepositUpdate {
	if t != nil {
		du.SetCreatedAt(*t)
	}
	return du
}

// SetUserID sets the "user" edge to the User entity by ID.
func (du *DepositUpdate) SetUserID(id uuid.UUID) *DepositUpdate {
	du.mutation.SetUserID(id)
	return du
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (du *DepositUpdate) SetNillableUserID(id *uuid.UUID) *DepositUpdate {
	if id != nil {
		du = du.SetUserID(*id)
	}
	return du
}

// SetUser sets the "user" edge to the User entity.
func (du *DepositUpdate) SetUser(u *User) *DepositUpdate {
	return du.SetUserID(u.ID)
}

// Mutation returns the DepositMutation object of the builder.
func (du *DepositUpdate) Mutation() *DepositMutation {
	return du.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (du *DepositUpdate) ClearUser() *DepositUpdate {
	du.mutation.ClearUser()
	return du
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (du *DepositUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, DepositMutation](ctx, du.sqlSave, du.mutation, du.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (du *DepositUpdate) SaveX(ctx context.Context) int {
	affected, err := du.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (du *DepositUpdate) Exec(ctx context.Context) error {
	_, err := du.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (du *DepositUpdate) ExecX(ctx context.Context) {
	if err := du.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (du *DepositUpdate) check() error {
	if v, ok := du.mutation.Status(); ok {
		if err := deposit.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Deposit.status": %w`, err)}
		}
	}
	return nil
}

func (du *DepositUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := du.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(deposit.Table, deposit.Columns, sqlgraph.NewFieldSpec(deposit.FieldID, field.TypeUUID))
	if ps := du.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := du.mutation.Amount(); ok {
		_spec.SetField(deposit.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := du.mutation.AddedAmount(); ok {
		_spec.AddField(deposit.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := du.mutation.Status(); ok {
		_spec.SetField(deposit.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := du.mutation.CreatedAt(); ok {
		_spec.SetField(deposit.FieldCreatedAt, field.TypeTime, value)
	}
	if du.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deposit.UserTable,
			Columns: []string{deposit.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := du.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deposit.UserTable,
			Columns: []string{deposit.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, du.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{deposit.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	du.mutation.done = true
	return n, nil
}

// DepositUpdateOne is the builder for updating a single Deposit entity.
type DepositUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *DepositMutation
}

// SetAmount sets the "amount" field.
func (duo *DepositUpdateOne) SetAmount(f float64) *DepositUpdateOne {
	duo.mutation.ResetAmount()
	duo.mutation.SetAmount(f)
	return duo
}

// AddAmount adds f to the "amount" field.
func (duo *DepositUpdateOne) AddAmount(f float64) *DepositUpdateOne {
	duo.mutation.AddAmount(f)
	return duo
}

// SetStatus sets the "status" field.
func (duo *DepositUpdateOne) SetStatus(d deposit.Status) *DepositUpdateOne {
	duo.mutation.SetStatus(d)
	return duo
}

// SetCreatedAt sets the "created_at" field.
func (duo *DepositUpdateOne) SetCreatedAt(t time.Time) *DepositUpdateOne {
	duo.mutation.SetCreatedAt(t)
	return duo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (duo *DepositUpdateOne) SetNillableCreatedAt(t *time.Time) *DepositUpdateOne {
	if t != nil {
		duo.SetCreatedAt(*t)
	}
	return duo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (duo *DepositUpdateOne) SetUserID(id uuid.UUID) *DepositUpdateOne {
	duo.mutation.SetUserID(id)
	return duo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (duo *DepositUpdateOne) SetNillableUserID(id *uuid.UUID) *DepositUpdateOne {
	if id != nil {
		duo = duo.SetUserID(*id)
	}
	return duo
}

// SetUser sets the "user" edge to the User entity.
func (duo *DepositUpdateOne) SetUser(u *User) *DepositUpdateOne {
	return duo.SetUserID(u.ID)
}

// Mutation returns the DepositMutation object of the builder.
func (duo *DepositUpdateOne) Mutation() *DepositMutation {
	return duo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (duo *DepositUpdateOne) ClearUser() *DepositUpdateOne {
	duo.mutation.ClearUser()
	return duo
}

// Where appends a list predicates to the DepositUpdate builder.
func (duo *DepositUpdateOne) Where(ps ...predicate.Deposit) *DepositUpdateOne {
	duo.mutation.Where(ps...)
	return duo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (duo *DepositUpdateOne) Select(field string, fields ...string) *DepositUpdateOne {
	duo.fields = append([]string{field}, fields...)
	return duo
}

// Save executes the query and returns the updated Deposit entity.
func (duo *DepositUpdateOne) Save(ctx context.Context) (*Deposit, error) {
	return withHooks[*Deposit, DepositMutation](ctx, duo.sqlSave, duo.mutation, duo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (duo *DepositUpdateOne) SaveX(ctx context.Context) *Deposit {
	node, err := duo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (duo *DepositUpdateOne) Exec(ctx context.Context) error {
	_, err := duo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (duo *DepositUpdateOne) ExecX(ctx context.Context) {
	if err := duo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (duo *DepositUpdateOne) check() error {
	if v, ok := duo.mutation.Status(); ok {
		if err := deposit.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Deposit.status": %w`, err)}
		}
	}
	return nil
}

func (duo *DepositUpdateOne) sqlSave(ctx context.Context) (_node *Deposit, err error) {
	if err := duo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(deposit.Table, deposit.Columns, sqlgraph.NewFieldSpec(deposit.FieldID, field.TypeUUID))
	id, ok := duo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Deposit.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := duo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, deposit.FieldID)
		for _, f := range fields {
			if !deposit.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != deposit.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := duo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := duo.mutation.Amount(); ok {
		_spec.SetField(deposit.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := duo.mutation.AddedAmount(); ok {
		_spec.AddField(deposit.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := duo.mutation.Status(); ok {
		_spec.SetField(deposit.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := duo.mutation.CreatedAt(); ok {
		_spec.SetField(deposit.FieldCreatedAt, field.TypeTime, value)
	}
	if duo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deposit.UserTable,
			Columns: []string{deposit.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := duo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deposit.UserTable,
			Columns: []string{deposit.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Deposit{config: duo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, duo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{deposit.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	duo.mutation.done = true
	return _node, nil
}
