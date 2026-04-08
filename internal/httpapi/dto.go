package httpapi

import "test-task/internal/greenapi"

type credentialsRequest struct {
	IDInstance string `json:"idInstance"`
	APIToken   string `json:"apiTokenInstance"`
}

type settingsRequest struct {
	credentialsRequest
}

type stateRequest struct {
	credentialsRequest
}

type sendMessagePageRequest struct {
	credentialsRequest
	greenapi.SendMessageRequest
}

type sendFilePageRequest struct {
	credentialsRequest
	greenapi.SendFileByURLRequest
}
