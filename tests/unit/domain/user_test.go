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
			"create valid user",
			args{
				fullName: "John Doe",
				password: "mypassword",
				email:    "john@example.com",
			},
			nil,
		},

		{
			"too short password",
			args{
				fullName: "John Doe",
				password: "123",
				email:    "john@example.com",
			},
			domain.ErrPasswordTooShort,
		},
		{
			"too short fullName",
			args{
				fullName: "JD",
				password: "mypassword",
				email:    "john@example.com",
			},
			domain.ErrTooShortFullName,
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
	hashedPassword, err := domain.HashPassword("secret123")
	assert.NoError(t, err, "should hash password without error")

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
			"valid password",
			args{
				oldPassword: hashedPassword,
				password:    "secret123",
			},
			nil,
		},
		{
			"invalid old password",
			args{
				oldPassword: "secret123", // this must be hashed
				password:    "secret123",
			},
			domain.ErrInternalServer,
		},
		{
			"invalid password",
			args{
				oldPassword: hashedPassword,
				password:    "wrongpassword",
			},
			domain.ErrInvalidCredentials,
		},
		{
			"short password",
			args{
				oldPassword: "secret123",
				password:    "short",
			},
			domain.ErrPasswordTooShort,
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
