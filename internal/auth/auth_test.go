package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestArgon2id(t *testing.T) {
	password := "pa$$word"
	hash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)

	match, _ := argon2id.ComparePasswordAndHash(password, hash)

	if !match {
		t.Errorf("password did not match hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWTFunctions(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	t.Run("Valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT failed: %v", err)
		}

		id, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("ValidateJWT failed: %v", err)
		}

		if id != userID {
			t.Errorf("Expected %v, got %v", userID, id)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Minute)
		if err != nil {
			t.Errorf("MakeJWT failed: %v", err)
		}

		id, err := ValidateJWT(token, secret)
		if err == nil {
			t.Errorf("Expected error for expired token, got none. Returned ID: %v", id)
		}
	})

	t.Run("Wrong secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT failed: %v", err)
		}
		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Errorf("Expected error for wrong secret, got none")
		}
	})
}
