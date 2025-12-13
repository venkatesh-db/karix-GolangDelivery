package domain

import (
	"strings"
	"time"
)

// CustomerID represents the globally unique identifier for a customer aggregate.
type CustomerID string

// Customer models the core attributes HelraChar requires for onboarding.
type Customer struct {
	id        CustomerID
	fullName  string
	email     string
	pan       string
	createdAt time.Time
}

// NewCustomer enforces invariant checks before a customer aggregate is created.
func NewCustomer(id CustomerID, fullName, email, pan string, createdAt time.Time) (*Customer, error) {
	errs := NewValidationError()
	name := strings.TrimSpace(fullName)
	mail := strings.TrimSpace(email)
	panValue := strings.TrimSpace(pan)

	if id == "" {
		errs.Add("id", "cannot be empty")
	}
	if name == "" {
		errs.Add("fullName", "is required")
	}
	if !strings.Contains(mail, "@") {
		errs.Add("email", "must contain '@'")
	}
	if len(panValue) != 10 {
		errs.Add("pan", "must be 10 characters")
	}
	if errs.HasErrors() {
		return nil, errs
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	return &Customer{
		id:        id,
		fullName:  name,
		email:     strings.ToLower(mail),
		pan:       strings.ToUpper(panValue),
		createdAt: createdAt,
	}, nil
}

// ID exposes the aggregate identifier.
func (c *Customer) ID() CustomerID {
	return c.id
}

// FullName returns the immutable full name captured at onboarding.
func (c *Customer) FullName() string {
	return c.fullName
}

// Email returns the normalized email address.
func (c *Customer) Email() string {
	return c.email
}

// PAN returns the 10 character Permanent Account Number reference.
func (c *Customer) PAN() string {
	return c.pan
}

// CreatedAt returns the timestamp when the aggregate was persisted.
func (c *Customer) CreatedAt() time.Time {
	return c.createdAt
}
