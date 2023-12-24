package Discord

import "time"

type Message struct {
	ID              string      `json:"id"`
	ChannelID       string      `json:"channel_id"`
	Author          DiscordUser `json:"author"`
	Content         string      `json:"content"`
	Timestamp       time.Time   `json:"timestamp"`
	EditedTimestamp *time.Time  `json:"edited_timestamp"`
	TTS             bool        `json:"tts"`
	MentionEveryone bool        `json:"mention_everyone"`
	Mentions        []struct{}  `json:"mentions"`
	MentionRoles    []string    `json:"mention_roles"`
	Attachments     []struct {
		ID                 string `json:"id"`
		Filename           string `json:"filename"`
		Size               int    `json:"size"`
		URL                string `json:"url"`
		ProxyURL           string `json:"proxy_url"`
		Height             int    `json:"height,omitempty"`
		Width              int    `json:"width,omitempty"`
		ContentScanVersion int    `json:"content_scan_version,omitempty"`
	} `json:"attachments"`
	Embeds            []interface{} `json:"embeds"`
	Reactions         []interface{} `json:"reactions"`
	Pinned            bool          `json:"pinned"`
	WebhookID         *string       `json:"webhook_id,omitempty"`
	Type              int           `json:"type"`
	Activity          interface{}   `json:"activity,omitempty"`
	Application       interface{}   `json:"application,omitempty"`
	ApplicationID     *string       `json:"application_id,omitempty"`
	MessageReference  interface{}   `json:"message_reference,omitempty"`
	Flags             int           `json:"flags"`
	ReferencedMessage *Message      `json:"referenced_message,omitempty"`
	Interaction       interface{}   `json:"interaction,omitempty"`
	Thread            interface{}   `json:"thread,omitempty"`
	Components        []interface{} `json:"components"`
	StickerItems      []interface{} `json:"sticker_items,omitempty"`
	Stickers          []interface{} `json:"stickers,omitempty"`
	Position          int           `json:"position,omitempty"`
	RoleSubscription  interface{}   `json:"role_subscription_data,omitempty"`
	Resolved          interface{}   `json:"resolved,omitempty"`
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Flags         int    `json:"flags"`
	// ... add other fields as needed
}
