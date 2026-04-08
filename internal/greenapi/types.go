package greenapi

type Credentials struct {
	IDInstance string `json:"idInstance"`
	APIToken   string `json:"apiTokenInstance"`
}

type SettingsResponse struct {
	WID                               string `json:"wid"`
	WebhookURL                        string `json:"webhookUrl"`
	WebhookURLToken                   string `json:"webhookUrlToken"`
	DelaySendMessagesMilliseconds     int    `json:"delaySendMessagesMilliseconds"`
	MarkIncomingMessagesReaded        string `json:"markIncomingMessagesReaded"`
	MarkIncomingMessagesReadedOnReply string `json:"markIncomingMessagesReadedOnReply"`
	OutgoingWebhook                   string `json:"outgoingWebhook"`
	OutgoingMessageWebhook            string `json:"outgoingMessageWebhook"`
	OutgoingAPIMessageWebhook         string `json:"outgoingAPIMessageWebhook"`
	IncomingWebhook                   string `json:"incomingWebhook"`
	DeviceWebhook                     string `json:"deviceWebhook"`
	StateWebhook                      string `json:"stateWebhook"`
	KeepOnlineStatus                  string `json:"keepOnlineStatus"`
	PollMessageWebhook                string `json:"pollMessageWebhook"`
	IncomingBlockWebhook              string `json:"incomingBlockWebhook"`
	IncomingCallWebhook               string `json:"incomingCallWebhook"`
	EditedMessageWebhook              string `json:"editedMessageWebhook"`
	DeletedMessageWebhook             string `json:"deletedMessageWebhook"`
}

type StateResponse struct {
	StateInstance string `json:"stateInstance"`
}

type SendMessageRequest struct {
	ChatID          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageID string `json:"quotedMessageId,omitempty"`
	LinkPreview     *bool  `json:"linkPreview,omitempty"`
	TypingTime      *int   `json:"typingTime,omitempty"`
}

type SendFileByURLRequest struct {
	ChatID          string `json:"chatId"`
	URLFile         string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption,omitempty"`
	QuotedMessageID string `json:"quotedMessageId,omitempty"`
	TypingTime      *int   `json:"typingTime,omitempty"`
	TypingType      string `json:"typingType,omitempty"`
}

type SendMessageResponse struct {
	IDMessage string `json:"idMessage"`
}
