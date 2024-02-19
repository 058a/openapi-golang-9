package item_test

import (
	"openapi/internal/domain/stock/item"
	"testing"
)

func TestNewName(t *testing.T) {
	// Setup
	t.Parallel()

	type args struct {
		v string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		errType error
	}{
		{
			name: "success",
			args: args{
				v: "test",
			},
			want:    "test",
			wantErr: false,
			errType: nil,
		},
		{
			name: "fail",
			args: args{
				v: "",
			},
			want:    "",
			wantErr: true,
			errType: item.ErrNameEmpty,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := item.NewName(tt.args.v)

			// Then
			if !tt.wantErr {
				// 正常系
				if err != nil {
					t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
					return						
				}
				if got.String() != tt.want {
					t.Errorf("NewName() = %v, want %v", got.String(), tt.want)
				}	
			}

			// 異常系
			if err != tt.errType {
				t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
