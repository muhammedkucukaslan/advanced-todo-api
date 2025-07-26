package testdomain

import (
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	type args struct {
		fullName string
		password string
		email    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// Validator library checks if email is valid, so we can skip that here
		{
			name: "should create valid user",
			args: args{
				fullName: "John Doe",
				password: "mypassword",
				email:    "john@example.com",
			},
			wantErr: nil,
		},

		{
			name: "should fail with weak password",
			args: args{
				fullName: "John Doe",
				password: "123",
				email:    "john@example.com",
			},
			wantErr: domain.ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewUser(tt.args.fullName, tt.args.password, tt.args.email)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assert.NotEmpty(t, got.Id)
				_, err := got.Id.MarshalBinary()
				assert.NoError(t, err, "Id should be valid UUID format")

				assert.NotNil(t, got.Password)
				assert.Greater(t, len(got.Password), 10, "hashed password should be longer than original")

				assert.Equal(t, tt.args.fullName, got.FullName)
				assert.Equal(t, tt.args.email, got.Email)
				assert.Equal(t, "USER", got.Role)

				assert.NotEqual(t, tt.args.password, got.Password)
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	// This is a hashed password for `secret123`
	hashedPassword := "$2y$10$IkBM5kx5MuMxCGLyHEepveo3GfcDHnR2H22wUjWSIitlp4DDVGwIu"

	type args struct {
		oldPassword string
		password    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid password",
			args: args{
				oldPassword: hashedPassword,
				password:    "secret123",
			},
			wantErr: nil,
		},
		{
			name: "invalid old password",
			args: args{
				oldPassword: "secret123", // this must be hashed
				password:    "secret123",
			},
			wantErr: domain.ErrInternalServer,
		},
		{
			name: "invalid password",
			args: args{
				oldPassword: hashedPassword,
				password:    "wrongpassword",
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "short password",
			args: args{
				oldPassword: "secret123",
				password:    "short",
			},
			wantErr: domain.ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &domain.User{
				Password: tt.args.oldPassword,
			}

			err := u.ValidatePassword(tt.args.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
