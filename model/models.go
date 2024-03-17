// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Balance     pgtype.Numeric   `json:"balance"`
	UserID      uuid.UUID        `json:"user_id"`
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type Category struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Icon        *string          `json:"icon"`
	Color       *string          `json:"color"`
	CatType     string           `json:"cat_type"`
	UserID      uuid.UUID        `json:"user_id"`
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type Token struct {
	ID         int32            `json:"id"`
	Token      string           `json:"token"`
	Identifier *string          `json:"identifier"`
	Mobile     bool             `json:"mobile"`
	UserID     uuid.UUID        `json:"user_id"`
	ExpiresAt  pgtype.Timestamp `json:"expires_at"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}

type Transaction struct {
	ID          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Note        *string          `json:"note"`
	Amount      pgtype.Numeric   `json:"amount"`
	Paid        bool             `json:"paid"`
	TType       string           `json:"t_type"`
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	UserID      uuid.UUID        `json:"user_id"`
	CategoryID  uuid.UUID        `json:"category_id"`
	AccountID   uuid.UUID        `json:"account_id"`
	HandledAt   pgtype.Timestamp `json:"handled_at"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type User struct {
	ID                uuid.UUID        `json:"id"`
	FirstName         string           `json:"first_name"`
	LastName          string           `json:"last_name"`
	Email             string           `json:"email"`
	Password          string           `json:"password"`
	Verified          bool             `json:"verified"`
	VerificationToken *string          `json:"verification_token"`
	Avatar            string           `json:"avatar"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type Workspace struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Currency    string           `json:"currency"`
	Language    string           `json:"language"`
	UserID      uuid.UUID        `json:"user_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}

type WorkspacesUser struct {
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	UserID      uuid.UUID        `json:"user_id"`
	Role        string           `json:"role"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}
