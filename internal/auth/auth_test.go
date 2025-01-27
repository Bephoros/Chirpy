package auth

import (
	"net/http"
	"testing"

	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("expected hashed password, got empty string")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword"

	// Generate hash
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test correct password
	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("expected passwords to match, got error: %v", err)
	}

	// Test incorrect password
	err = CheckPasswordHash("wrongpassword", hashedPassword)
	if err == nil {
		t.Errorf("expected error for incorrect password, got nil")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectedError string
	}{
		{
			name: "Token válido",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			expectedToken: "valid_token",
			expectedError: "",
		},
		{
			name: "Sin encabezado Authorization",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedToken: "",
			expectedError: "authorization header is missing",
		},
		{
			name: "Encabezado sin prefijo Bearer",
			headers: http.Header{
				"Authorization": []string{"InvalidTokenString"},
			},
			expectedToken: "",
			expectedError: "authorization header is not a bearer token",
		},
		{
			name: "Token vacío después del prefijo",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			expectedToken: "",
			expectedError: "bearer token is empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.headers)

			if tc.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tc.expectedError)
				} else if err.Error() != tc.expectedError {
					t.Errorf("expected error '%s', got '%s'", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("expected token '%s', got '%s'", tc.expectedToken, token)
				}
			}
		})
	}
}
