package filemaker

import (
	"errors"
	"net/http"
	"testing"
)

func TestFileMakerError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *FileMakerError
		expected string
	}{
		{
			name: "full error with code and message",
			err: &FileMakerError{
				Code:       "401",
				Message:    "Invalid credentials",
				HTTPStatus: http.StatusUnauthorized,
			},
			expected: "FileMaker error 401: Invalid credentials (HTTP 401)",
		},
		{
			name: "error without code",
			err: &FileMakerError{
				Message:    "Server error",
				HTTPStatus: http.StatusInternalServerError,
			},
			expected: "FileMaker error: Server error (HTTP 500)",
		},
		{
			name: "error with only status",
			err: &FileMakerError{
				HTTPStatus: http.StatusBadRequest,
			},
			expected: "FileMaker error (HTTP 400)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFileMakerError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &FileMakerError{
		Code:       "500",
		Message:    "Server error",
		HTTPStatus: http.StatusInternalServerError,
		Err:        underlyingErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "error with field",
			err: &ValidationError{
				Field:   "username",
				Message: "cannot be empty",
			},
			expected: "validation error on field 'username': cannot be empty",
		},
		{
			name: "error without field",
			err: &ValidationError{
				Message: "invalid input",
			},
			expected: "validation error: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &ValidationError{
		Field:   "test",
		Message: "test error",
		Err:     underlyingErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestAuthenticationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AuthenticationError
		expected string
	}{
		{
			name: "with message",
			err: &AuthenticationError{
				Message: "invalid credentials",
			},
			expected: "authentication error: invalid credentials",
		},
		{
			name:     "without message",
			err:      &AuthenticationError{},
			expected: "authentication error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAuthenticationError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &AuthenticationError{
		Message: "test",
		Err:     underlyingErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestNetworkError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *NetworkError
		expected string
	}{
		{
			name: "with message",
			err: &NetworkError{
				Message: "connection refused",
			},
			expected: "network error: connection refused",
		},
		{
			name:     "without message",
			err:      &NetworkError{},
			expected: "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNetworkError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &NetworkError{
		Message: "test",
		Err:     underlyingErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestTimeoutError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *TimeoutError
		expected string
	}{
		{
			name: "with message",
			err: &TimeoutError{
				Message: "request timeout",
			},
			expected: "timeout error: request timeout",
		},
		{
			name:     "without message",
			err:      &TimeoutError{},
			expected: "timeout error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTimeoutError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &TimeoutError{
		Message: "test",
		Err:     underlyingErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestIsFileMakerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is filemaker error",
			err:  &FileMakerError{Code: "401"},
			want: true,
		},
		{
			name: "wrapped filemaker error",
			err:  errors.Join(&FileMakerError{Code: "401"}, errors.New("other")),
			want: true,
		},
		{
			name: "not filemaker error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileMakerError(tt.err); got != tt.want {
				t.Errorf("IsFileMakerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is validation error",
			err:  &ValidationError{Field: "test"},
			want: true,
		},
		{
			name: "not validation error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is authentication error",
			err:  &AuthenticationError{Message: "test"},
			want: true,
		},
		{
			name: "not authentication error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthenticationError(tt.err); got != tt.want {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is network error",
			err:  &NetworkError{Message: "test"},
			want: true,
		},
		{
			name: "not network error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetworkError(tt.err); got != tt.want {
				t.Errorf("IsNetworkError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is timeout error",
			err:  &TimeoutError{Message: "test"},
			want: true,
		},
		{
			name: "not timeout error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTimeoutError(tt.err); got != tt.want {
				t.Errorf("IsTimeoutError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "host unavailable error",
			err:  &FileMakerError{Code: "952"},
			want: true,
		},
		{
			name: "too many files open error",
			err:  &FileMakerError{Code: "953"},
			want: true,
		},
		{
			name: "too many requests",
			err:  &FileMakerError{HTTPStatus: http.StatusTooManyRequests},
			want: true,
		},
		{
			name: "internal server error",
			err:  &FileMakerError{HTTPStatus: http.StatusInternalServerError},
			want: true,
		},
		{
			name: "bad gateway",
			err:  &FileMakerError{HTTPStatus: http.StatusBadGateway},
			want: true,
		},
		{
			name: "service unavailable",
			err:  &FileMakerError{HTTPStatus: http.StatusServiceUnavailable},
			want: true,
		},
		{
			name: "gateway timeout",
			err:  &FileMakerError{HTTPStatus: http.StatusGatewayTimeout},
			want: true,
		},
		{
			name: "network error",
			err:  &NetworkError{Message: "connection refused"},
			want: true,
		},
		{
			name: "timeout error",
			err:  &TimeoutError{Message: "timeout"},
			want: true,
		},
		{
			name: "non-retryable filemaker error",
			err:  &FileMakerError{Code: "401", HTTPStatus: http.StatusUnauthorized},
			want: false,
		},
		{
			name: "non-retryable generic error",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "authentication error",
			err:  &AuthenticationError{Message: "invalid credentials"},
			want: true,
		},
		{
			name: "filemaker invalid username",
			err:  &FileMakerError{Code: "212"},
			want: true,
		},
		{
			name: "filemaker no access privileges",
			err:  &FileMakerError{Code: "214"},
			want: true,
		},
		{
			name: "filemaker invalid token",
			err:  &FileMakerError{Code: "952"},
			want: true,
		},
		{
			name: "http unauthorized",
			err:  &FileMakerError{HTTPStatus: http.StatusUnauthorized},
			want: true,
		},
		{
			name: "http forbidden",
			err:  &FileMakerError{HTTPStatus: http.StatusForbidden},
			want: true,
		},
		{
			name: "non-auth error",
			err:  &FileMakerError{Code: "500"},
			want: false,
		},
		{
			name: "generic error",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthError(tt.err); got != tt.want {
				t.Errorf("IsAuthError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFileMakerError(t *testing.T) {
	tests := []struct {
		name       string
		response   *ResponseData
		httpStatus int
		wantErr    bool
		wantCode   string
		wantMsg    string
	}{
		{
			name:       "nil response",
			response:   nil,
			httpStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantMsg:    "no response data",
		},
		{
			name: "empty messages",
			response: &ResponseData{
				Messages: []Message{},
			},
			httpStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantMsg:    "unknown error",
		},
		{
			name: "success code",
			response: &ResponseData{
				Messages: []Message{{Code: "0", Message: "OK"}},
			},
			httpStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "error code",
			response: &ResponseData{
				Messages: []Message{{Code: "401", Message: "Invalid credentials"}},
			},
			httpStatus: http.StatusUnauthorized,
			wantErr:    true,
			wantCode:   "401",
			wantMsg:    "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseFileMakerError(tt.response, tt.httpStatus)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFileMakerError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				fmErr, ok := err.(*FileMakerError)
				if !ok {
					t.Errorf("ParseFileMakerError() returned non-FileMakerError: %T", err)
					return
				}

				if tt.wantCode != "" && fmErr.Code != tt.wantCode {
					t.Errorf("ParseFileMakerError() Code = %v, want %v", fmErr.Code, tt.wantCode)
				}

				if tt.wantMsg != "" && fmErr.Message != tt.wantMsg {
					t.Errorf("ParseFileMakerError() Message = %v, want %v", fmErr.Message, tt.wantMsg)
				}
			}
		})
	}
}
