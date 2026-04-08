package httpapi

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"test-task/internal/greenapi"
)

const maxMessageLength = 20000

func ValidateSendMessage(req greenapi.SendMessageRequest) error {
	if err := validateChatID(req.ChatID); err != nil {
		return err
	}
	if strings.TrimSpace(req.Message) == "" {
		return errors.New("message is required")
	}
	if len([]rune(req.Message)) > maxMessageLength {
		return fmt.Errorf("message must be <= %d characters", maxMessageLength)
	}
	return validateTypingTime(req.TypingTime)
}

func ValidateSendFileByURL(req greenapi.SendFileByURLRequest) error {
	if err := validateChatID(req.ChatID); err != nil {
		return err
	}
	if strings.TrimSpace(req.URLFile) == "" {
		return errors.New("urlFile is required")
	}
	parsedURL, err := url.Parse(req.URLFile)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.Host == "" {
		return errors.New("urlFile must be a valid https URL")
	}
	if strings.TrimSpace(req.FileName) == "" {
		return errors.New("fileName is required")
	}
	if ext := filepath.Ext(req.FileName); ext == "" {
		return errors.New("fileName must include an extension")
	}
	if len([]rune(req.Caption)) > maxMessageLength {
		return fmt.Errorf("caption must be <= %d characters", maxMessageLength)
	}
	if req.TypingType != "" && req.TypingType != "recording" {
		return errors.New(`typingType must be "recording" when provided`)
	}
	return validateTypingTime(req.TypingTime)
}

func validateChatID(chatID string) error {
	chatID = strings.TrimSpace(chatID)
	if chatID == "" {
		return errors.New("chatId is required")
	}
	if !strings.HasSuffix(chatID, "@c.us") && !strings.HasSuffix(chatID, "@g.us") {
		return errors.New(`chatId must end with "@c.us" or "@g.us"`)
	}
	return nil
}

func validateTypingTime(typingTime *int) error {
	if typingTime == nil {
		return nil
	}
	if *typingTime < 1000 || *typingTime > 20000 {
		return errors.New("typingTime must be between 1000 and 20000 milliseconds")
	}
	return nil
}
