package request

import "github.com/google/uuid"

type RegisterRequest struct {
	Fullname string `json:"fullname" validate:"required,max=128"`
	Username string `json:"username" validate:"required,min=2,max=64"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Username string `json:"username" validate:"omitempty,min=2,max=64"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type CreateCommentRequest struct {
	Message string `json:"message" validate:"required,min=1,max=5000"`
}

type CreateAIJobRequest struct {
	FileID uuid.UUID `json:"file_id" validate:"required"`
	Type   string    `json:"type" validate:"required,oneof=TRANSLATE SUMMARIZE COMMENT"`
	Input  string    `json:"input" validate:"omitempty,max=10000"`
}
