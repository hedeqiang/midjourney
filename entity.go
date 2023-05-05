package midjourney

type Payload struct {
	Type          int    `json:"type"`
	ApplicationID string `json:"application_id"`
	GuildID       string `json:"guild_id"`
	ChannelID     string `json:"channel_id"`
	SessionID     string `json:"session_id"`
	Data          Data   `json:"data"`
}
type Options struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ApplicationCommandOptions struct {
	Type        int    `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}
type ApplicationCommand struct {
	ID                       string                      `json:"id"`
	ApplicationID            string                      `json:"application_id"`
	Version                  string                      `json:"version"`
	DefaultMemberPermissions interface{}                 `json:"default_member_permissions"`
	Type                     int                         `json:"type"`
	Nsfw                     bool                        `json:"nsfw"`
	Name                     string                      `json:"name"`
	Description              string                      `json:"description"`
	DmPermission             bool                        `json:"dm_permission"`
	Options                  []ApplicationCommandOptions `json:"options"`
}
type Data struct {
	Version            string             `json:"version"`
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Type               int                `json:"type"`
	Options            []Options          `json:"options"`
	ApplicationCommand ApplicationCommand `json:"application_command"`
	Attachments        []interface{}      `json:"attachments"`
}

type UpscalePayload struct {
	Type          int         `json:"type"`
	Nonce         string      `json:"nonce"`
	GuildId       string      `json:"guild_id"`
	ChannelId     string      `json:"channel_id"`
	MessageFlags  int         `json:"message_flags"`
	MessageId     string      `json:"message_id"`
	ApplicationId string      `json:"application_id"`
	SessionId     string      `json:"session_id"`
	Data          UpscaleData `json:"data"`
}

type UpscaleData struct {
	ComponentType int    `json:"component_type"`
	CustomId      string `json:"custom_id"`
}

type Attachments struct {
	Attachments []Attachment `json:"attachments"`
}
type Attachment struct {
	Id             int    `json:"id"`
	UploadUrl      string `json:"upload_url"`
	UploadFilename string `json:"upload_filename"`
}

type Files struct {
	Files []File `json:"files"`
}
type File struct {
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
	Id       int    `json:"id"`
}
