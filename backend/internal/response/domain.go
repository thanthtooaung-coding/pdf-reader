package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	RoleID    uuid.UUID `json:"role_id"`
	RoleCode  string    `json:"role_code,omitempty"`
	Fullname  string    `json:"fullname,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsEnable  bool      `json:"is_enable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type OTPSentResponse struct {
	Message          string `json:"message"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
	OTP              string `json:"otp,omitempty"`
}

type WorkspaceResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	FileCount int64     `json:"file_count,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	WorkspaceID      uuid.UUID `json:"workspace_id"`
	OriginalFileName string    `json:"original_file_name"`
	Name             string    `json:"name"`
	URL              string    `json:"url"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CommentResponse struct {
	ID        uuid.UUID `json:"id"`
	FileID    uuid.UUID `json:"file_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type AIJobResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	FileID      uuid.UUID `json:"file_id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Input       string    `json:"input,omitempty"`
	Output      string    `json:"output,omitempty"`
	DurationMS  int64     `json:"duration_ms"`
	RetryCount  int       `json:"retry_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
