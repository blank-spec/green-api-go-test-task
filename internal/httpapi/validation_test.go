package httpapi

import (
	"strings"
	"testing"

	"test-task/internal/greenapi"
)

func TestValidateSendMessage(t *testing.T) {
	typingTime := 3000

	tests := []struct {
		name    string
		req     greenapi.SendMessageRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: greenapi.SendMessageRequest{
				ChatID:     "79991234567@c.us",
				Message:    "hello",
				TypingTime: &typingTime,
			},
		},
		{
			name: "invalid chat id",
			req: greenapi.SendMessageRequest{
				ChatID:  "79991234567",
				Message: "hello",
			},
			wantErr: true,
		},
		{
			name: "empty message",
			req: greenapi.SendMessageRequest{
				ChatID:  "79991234567@c.us",
				Message: "   ",
			},
			wantErr: true,
		},
		{
			name: "too long message",
			req: greenapi.SendMessageRequest{
				ChatID:  "79991234567@c.us",
				Message: strings.Repeat("a", maxMessageLength+1),
			},
			wantErr: true,
		},
		{
			name: "typing time out of range",
			req: greenapi.SendMessageRequest{
				ChatID:     "79991234567@c.us",
				Message:    "hello",
				TypingTime: intPtr(999),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendMessage(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSendFileByURL(t *testing.T) {
	tests := []struct {
		name    string
		req     greenapi.SendFileByURLRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: greenapi.SendFileByURLRequest{
				ChatID:   "1234567890@g.us",
				URLFile:  "https://example.com/report.pdf",
				FileName: "report.pdf",
			},
		},
		{
			name: "invalid url scheme",
			req: greenapi.SendFileByURLRequest{
				ChatID:   "79991234567@c.us",
				URLFile:  "http://example.com/report.pdf",
				FileName: "report.pdf",
			},
			wantErr: true,
		},
		{
			name: "missing extension",
			req: greenapi.SendFileByURLRequest{
				ChatID:   "79991234567@c.us",
				URLFile:  "https://example.com/report",
				FileName: "report",
			},
			wantErr: true,
		},
		{
			name: "invalid typing type",
			req: greenapi.SendFileByURLRequest{
				ChatID:     "79991234567@c.us",
				URLFile:    "https://example.com/report.pdf",
				FileName:   "report.pdf",
				TypingType: "typing",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendFileByURL(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSendFileByURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func intPtr(v int) *int {
	return &v
}
