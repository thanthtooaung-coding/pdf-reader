package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/app"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

var testApp *app.Application

func TestMain(m *testing.M) {
	if os.Getenv("DATABASE_URL") == "" {
		fmt.Println("skip integration tests: DATABASE_URL not set")
		os.Exit(0)
	}

	cfg := config.Load()
	cfg.LogLevel = "error"
	cfg.DBAutoMigrate = true
	cfg.OTPDebugReturn = true
	cfg.SeedDefaultAdmin = true
	cfg.JWTSecret = "integration-test-secret"

	storageDir, err := os.MkdirTemp("", "pdfreader-api-test-*")
	if err != nil {
		fmt.Printf("setup storage: %v\n", err)
		os.Exit(1)
	}
	cfg.StoragePath = storageDir

	application, err := app.Bootstrap(cfg)
	if err != nil {
		fmt.Printf("bootstrap: %v\n", err)
		_ = os.RemoveAll(storageDir)
		os.Exit(1)
	}
	testApp = application

	code := m.Run()

	_ = testApp.Close()
	_ = os.RemoveAll(storageDir)
	os.Exit(code)
}

type apiEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func doJSON(t *testing.T, method, path string, body any, token string) (int, apiEnvelope) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := testApp.App.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	var env apiEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil && resp.StatusCode != http.StatusOK {
		t.Fatalf("decode response: %v", err)
	}
	return resp.StatusCode, env
}

func decodeData[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	return out
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%s@pdfreader.test", prefix, uuid.NewString()[:8])
}

func registerUser(t *testing.T, fullname, username, email, password string) (string, response.AuthResponse) {
	t.Helper()

	status, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/request", map[string]string{
		"fullname": fullname,
		"username": username,
		"email":    email,
		"password": password,
	}, "")
	if status != http.StatusOK || !env.Success {
		t.Fatalf("register request: status=%d env=%+v", status, env)
	}

	otpResp := decodeData[response.OTPSentResponse](t, env.Data)
	if otpResp.OTP == "" {
		t.Fatal("expected otp in debug response")
	}

	status, env = doJSON(t, http.MethodPost, "/api/v1/auth/register/verify", map[string]string{
		"email": email,
		"otp":   otpResp.OTP,
	}, "")
	if status != http.StatusCreated || !env.Success {
		t.Fatalf("register verify: status=%d env=%+v", status, env)
	}

	auth := decodeData[response.AuthResponse](t, env.Data)
	if auth.AccessToken == "" {
		t.Fatal("expected access token")
	}
	return auth.AccessToken, auth
}

func uploadPDF(t *testing.T, token, workspaceID string) response.FileResponse {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "sample.pdf")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(minimalPDF); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	path := fmt.Sprintf("/api/v1/workspaces/%s/files", workspaceID)
	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := testApp.App.Test(req, -1)
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	defer resp.Body.Close()

	var env apiEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode upload: %v", err)
	}
	if resp.StatusCode != http.StatusCreated || !env.Success {
		t.Fatalf("upload pdf: status=%d env=%+v", resp.StatusCode, env)
	}
	return decodeData[response.FileResponse](t, env.Data)
}

var minimalPDF = []byte("%PDF-1.4\n1 0 obj<<>>endobj\ntrailer<<>>\n%%EOF\n")

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp, err := testApp.App.Test(req, -1)
	if err != nil {
		t.Fatalf("healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want 200", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status field=%q", body["status"])
	}
}

func TestAuthAPIs(t *testing.T) {
	email := uniqueEmail("auth")
	username := "authuser_" + uuid.NewString()[:8]
	password := "Password123!"

	t.Run("register request", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/request", map[string]string{
			"fullname": "Auth Test",
			"username": username,
			"email":    email,
			"password": password,
		}, "")
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		otpResp := decodeData[response.OTPSentResponse](t, env.Data)
		if otpResp.OTP == "" {
			t.Fatal("missing otp")
		}
	})

	t.Run("register resend", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/resend", map[string]string{
			"email": email,
		}, "")
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})

	var token string
	t.Run("register verify", func(t *testing.T) {
		_, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/resend", map[string]string{
			"email": email,
		}, "")
		otpResp := decodeData[response.OTPSentResponse](t, env.Data)

		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/verify", map[string]string{
			"email": email,
			"otp":   otpResp.OTP,
		}, "")
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		auth := decodeData[response.AuthResponse](t, env.Data)
		token = auth.AccessToken
	})

	t.Run("login request", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
			"email":    email,
			"username": username,
			"password": password,
		}, "")
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})

	t.Run("login verify", func(t *testing.T) {
		_, env := doJSON(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
			"email":    email,
			"username": username,
			"password": password,
		}, "")
		otpResp := decodeData[response.OTPSentResponse](t, env.Data)

		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/login/verify", map[string]string{
			"email": email,
			"otp":   otpResp.OTP,
		}, "")
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		auth := decodeData[response.AuthResponse](t, env.Data)
		if auth.AccessToken == "" {
			t.Fatal("missing token")
		}
		token = auth.AccessToken
	})

	t.Run("users me unauthorized", func(t *testing.T) {
		status, _ := doJSON(t, http.MethodGet, "/api/v1/users/me", nil, "")
		if status != http.StatusUnauthorized {
			t.Fatalf("status=%d want 401", status)
		}
	})

	t.Run("users me", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/users/me", nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		me := decodeData[response.UserResponse](t, env.Data)
		if me.Email != email || me.Username != username {
			t.Fatalf("user mismatch: %+v", me)
		}
	})
}

func TestWorkspaceFileCommentAIJobAPIs(t *testing.T) {
	email := uniqueEmail("flow")
	username := "flow_" + uuid.NewString()[:8]
	password := "Password123!"
	token, _ := registerUser(t, "Flow User", username, email, password)

	var workspaceID string
	t.Run("create workspace", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/workspaces", map[string]string{
			"name": "My PDF Workspace",
		}, token)
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		ws := decodeData[response.WorkspaceResponse](t, env.Data)
		workspaceID = ws.ID.String()
	})

	t.Run("list workspaces", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/workspaces?page=1&limit=10", nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		page := decodeData[response.PagedResponse[response.WorkspaceResponse]](t, env.Data)
		if len(page.Items) == 0 {
			t.Fatal("expected workspaces")
		}
	})

	t.Run("get workspace", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		ws := decodeData[response.WorkspaceResponse](t, env.Data)
		if ws.ID.String() != workspaceID {
			t.Fatalf("workspace id mismatch")
		}
	})

	var fileID string
	t.Run("upload file", func(t *testing.T) {
		file := uploadPDF(t, token, workspaceID)
		fileID = file.ID.String()
		if file.Type != "pdf" {
			t.Fatalf("type=%q want pdf", file.Type)
		}
	})

	t.Run("reject non pdf upload", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "notes.txt")
		_, _ = part.Write([]byte("hello"))
		_ = writer.Close()

		path := fmt.Sprintf("/api/v1/workspaces/%s/files", workspaceID)
		req := httptest.NewRequest(http.MethodPost, path, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := testApp.App.Test(req, -1)
		if err != nil {
			t.Fatalf("upload txt: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status=%d want 400", resp.StatusCode)
		}
	})

	t.Run("list workspace files", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/files", nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		var files []response.FileResponse
		if err := json.Unmarshal(env.Data, &files); err != nil {
			t.Fatalf("decode files: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("expected files")
		}
	})

	t.Run("get file", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/files/"+fileID, nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		file := decodeData[response.FileResponse](t, env.Data)
		if file.ID.String() != fileID {
			t.Fatalf("file id mismatch")
		}
	})

	t.Run("download file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/files/"+fileID+"/download", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := testApp.App.Test(req, -1)
		if err != nil {
			t.Fatalf("download: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status=%d want 200", resp.StatusCode)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(data, minimalPDF) {
			t.Fatal("downloaded content mismatch")
		}
	})

	t.Run("create comment", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/files/"+fileID+"/comments", map[string]string{
			"message": "Looks good on page 1",
		}, token)
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		comment := decodeData[response.CommentResponse](t, env.Data)
		if comment.Message != "Looks good on page 1" {
			t.Fatalf("comment=%q", comment.Message)
		}
	})

	t.Run("list comments", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/files/"+fileID+"/comments", nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		var comments []response.CommentResponse
		if err := json.Unmarshal(env.Data, &comments); err != nil {
			t.Fatalf("decode comments: %v", err)
		}
		if len(comments) == 0 {
			t.Fatal("expected comments")
		}
	})

	var jobID string
	t.Run("create ai job summarize", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/ai-jobs", map[string]any{
			"file_id": fileID,
			"type":    "SUMMARIZE",
		}, token)
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		job := decodeData[response.AIJobResponse](t, env.Data)
		jobID = job.ID.String()
	})

	t.Run("create ai job translate", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/ai-jobs", map[string]any{
			"file_id": fileID,
			"type":    "TRANSLATE",
			"input":   "my",
		}, token)
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})

	t.Run("create ai job comment", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/ai-jobs", map[string]any{
			"file_id": fileID,
			"type":    "COMMENT",
		}, token)
		if status != http.StatusCreated || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})

	t.Run("get ai job", func(t *testing.T) {
		deadline := time.Now().Add(3 * time.Second)
		for {
			status, env := doJSON(t, http.MethodGet, "/api/v1/ai-jobs/"+jobID, nil, token)
			if status != http.StatusOK || !env.Success {
				t.Fatalf("status=%d env=%+v", status, env)
			}
			job := decodeData[response.AIJobResponse](t, env.Data)
			if job.Status == "completed" {
				if job.Output == "" {
					t.Fatal("expected output")
				}
				return
			}
			if time.Now().After(deadline) {
				t.Fatalf("job not completed, status=%s", job.Status)
			}
			time.Sleep(100 * time.Millisecond)
		}
	})

	t.Run("list workspace ai jobs", func(t *testing.T) {
		status, env := doJSON(t, http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/ai-jobs", nil, token)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
		var jobs []response.AIJobResponse
		if err := json.Unmarshal(env.Data, &jobs); err != nil {
			t.Fatalf("decode jobs: %v", err)
		}
		if len(jobs) < 3 {
			t.Fatalf("expected at least 3 jobs, got %d", len(jobs))
		}
	})

	t.Run("forbidden workspace access", func(t *testing.T) {
		otherToken, _ := registerUser(t, "Other", "other_"+uuid.NewString()[:8], uniqueEmail("other"), password)
		status, _ := doJSON(t, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, otherToken)
		if status != http.StatusForbidden {
			t.Fatalf("status=%d want 403", status)
		}
	})
}

func TestProtectedRoutesRequireJWT(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/users/me"},
		{http.MethodPost, "/api/v1/workspaces"},
		{http.MethodGet, "/api/v1/workspaces"},
		{http.MethodGet, "/api/v1/files/" + uuid.NewString()},
		{http.MethodPost, "/api/v1/ai-jobs"},
	}

	for _, tc := range routes {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			status, _ := doJSON(t, tc.method, tc.path, map[string]string{"name": "x"}, "")
			if status != http.StatusUnauthorized {
				t.Fatalf("status=%d want 401", status)
			}
		})
	}
}

func TestAuthValidationErrors(t *testing.T) {
	t.Run("register missing fields", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/register/request", map[string]string{
			"email": "bad",
		}, "")
		if status != http.StatusBadRequest || env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})

	t.Run("login invalid credentials", func(t *testing.T) {
		status, env := doJSON(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
			"email":    uniqueEmail("nouser"),
			"username": "nobody",
			"password": "Password123!",
		}, "")
		if status != http.StatusUnauthorized || env.Success {
			t.Fatalf("status=%d env=%+v", status, env)
		}
	})
}
