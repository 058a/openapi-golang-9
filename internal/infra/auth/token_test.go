package auth

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestEncodeToken(t *testing.T) {
	t.Parallel()

	type args struct {
		userId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userId: uuid.New(),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := EncodeToken(tt.args.userId)

			// Then
			if !tt.wantErr {
				if err != nil {
					t.Errorf("EncodeToken() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got == "" {
					t.Errorf("EncodeToken() = %v, want %v", got, tt.want)
				}
				token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("unexpected signing method")
					}
					return []byte("secret"), nil
				})
				if err != nil {
					t.Errorf("EncodeToken() error = %v, wantErr %v", err, tt.wantErr)
				}

				if !token.Valid {
					t.Errorf("EncodeToken() = %v, want %v", got, tt.want)
				}
				return
			}

			if err == nil {
				t.Errorf("EncodeToken() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeToken(t *testing.T) {
	// Setup
	t.Parallel()

	userId := uuid.New()

	type args struct {
		userId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *jwtCustomClaims
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userId: userId,
			},
			want: &jwtCustomClaims{
				UserId: userId,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			token, err := EncodeToken(tt.args.userId)
			if err != nil {
				t.Fatal(err)
			}

			got, err := DecodeToken(token)
			if !tt.wantErr {
				if err != nil {
					t.Errorf("DecodeToken() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got.UserId != tt.want.UserId {
					t.Errorf("DecodeToken() = %v, want %v", got, tt.want)
				}
				return
			}

			if err == nil {
				t.Errorf("DecodeToken() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
