package response

// HealthResponse is returned by GET /healthz.
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"pdf-reader-backend"`
}

// ErrorEnvelope is the standard error response body.
type ErrorEnvelope struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"validation failed"`
}

type OTPEnvelope struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message,omitempty" example:"otp_sent"`
	Data    OTPSentResponse `json:"data"`
}

type AuthEnvelope struct {
	Success bool         `json:"success" example:"true"`
	Data    AuthResponse `json:"data"`
}

type UserEnvelope struct {
	Success bool         `json:"success" example:"true"`
	Data    UserResponse `json:"data"`
}

type WorkspaceEnvelope struct {
	Success bool              `json:"success" example:"true"`
	Data    WorkspaceResponse `json:"data"`
}

type WorkspacePage struct {
	Items []WorkspaceResponse `json:"items"`
	Meta  PageMeta            `json:"meta"`
}

type WorkspaceListEnvelope struct {
	Success bool          `json:"success" example:"true"`
	Data    WorkspacePage `json:"data"`
}

type FileEnvelope struct {
	Success bool         `json:"success" example:"true"`
	Data    FileResponse `json:"data"`
}

type FileListEnvelope struct {
	Success bool           `json:"success" example:"true"`
	Data    []FileResponse `json:"data"`
}

type CommentEnvelope struct {
	Success bool            `json:"success" example:"true"`
	Data    CommentResponse `json:"data"`
}

type CommentListEnvelope struct {
	Success bool              `json:"success" example:"true"`
	Data    []CommentResponse `json:"data"`
}

type AIJobEnvelope struct {
	Success bool          `json:"success" example:"true"`
	Data    AIJobResponse `json:"data"`
}

type AIJobListEnvelope struct {
	Success bool            `json:"success" example:"true"`
	Data    []AIJobResponse `json:"data"`
}
