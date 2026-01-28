package validator

import (
	"strings"
	"testing"
)

type TestUser struct {
	Name string `validate:"required,min=3"`
	Age  int    `validate:"gte=18"`
}

func TestValidator(t *testing.T) {
	tests := []struct {
		name    string
		input   TestUser
		wantErr bool
		errTag  string
	}{
		{
			name:    "valid user",
			input:   TestUser{Name: "John", Age: 18},
			wantErr: false,
		},
		{
			name:    "short name",
			input:   TestUser{Name: "Jo", Age: 20},
			wantErr: true,
			errTag:  "min",
		},
		{
			name:    "missing name",
			input:   TestUser{Age: 20},
			wantErr: true,
			errTag:  "required",
		},
		{
			name:    "underage",
			input:   TestUser{Name: "John", Age: 17},
			wantErr: true,
			errTag:  "gte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StructValidator.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				valErr, ok := err.(*ValidateError)
				if !ok {
					t.Errorf("Expected ValidateError, got %T", err)
				} else if valErr.Tag != tt.errTag {
					t.Errorf("Expected error tag %s, got %s", tt.errTag, valErr.Tag)
				}
				if !strings.Contains(err.Error(), "failed on "+tt.errTag) {
					t.Errorf("Error message %q does not contain tag %q", err.Error(), tt.errTag)
				}
			}
		})
	}
}
