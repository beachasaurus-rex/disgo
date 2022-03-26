package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

// WebhookType (https: //discord.com/developers/docs/resources/webhook#webhook-object-webhook-types)
type WebhookType int

// All WebhookType(s)
const (
	WebhookTypeIncoming WebhookType = iota + 1
	WebhookTypeChannelFollower
	WebhookTypeApplication
)

// Webhook (https://discord.com/developers/docs/resources/webhook) is a way to post messages to Discord using the Discord API which do not require bot authentication or use.
type Webhook interface {
	json.Marshaler
	Type() WebhookType
	ID() snowflake.Snowflake
	Name() string
	Avatar() *string
	AvatarURL(opts ...CDNOpt) *string
	webhook()
}

type UnmarshalWebhook struct {
	Webhook
}

func (w *UnmarshalWebhook) UnmarshalJSON(data []byte) error {
	var wType struct {
		Type WebhookType `json:"type"`
	}

	if err := json.Unmarshal(data, &wType); err != nil {
		return err
	}

	var (
		webhook Webhook
		err     error
	)

	switch wType.Type {
	case WebhookTypeIncoming:
		var v IncomingWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	case WebhookTypeChannelFollower:
		var v ChannelFollowerWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	case WebhookTypeApplication:
		var v ApplicationWebhook
		err = json.Unmarshal(data, &v)
		webhook = v

	default:
		err = fmt.Errorf("unkown webhook with type %d received", wType.Type)
	}

	if err != nil {
		return err
	}

	w.Webhook = webhook
	return nil
}

var _ Webhook = (*IncomingWebhook)(nil)

type IncomingWebhook struct {
	id            snowflake.Snowflake
	name          string
	avatar        *string
	ChannelID     snowflake.Snowflake  `json:"channel_id"`
	GuildID       snowflake.Snowflake  `json:"guild_id"`
	Token         string               `json:"token"`
	ApplicationID *snowflake.Snowflake `json:"application_id"`
	User          User                 `json:"user"`
}

func (w *IncomingWebhook) UnmarshalJSON(data []byte) error {
	type incomingWebhook IncomingWebhook
	var v struct {
		ID     snowflake.Snowflake `json:"id"`
		Name   string              `json:"name"`
		Avatar *string             `json:"avatar"`
		incomingWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = IncomingWebhook(v.incomingWebhook)
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	return nil
}

func (w IncomingWebhook) MarshalJSON() ([]byte, error) {
	type incomingWebhook IncomingWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		incomingWebhook
	}{
		Type:            w.Type(),
		incomingWebhook: incomingWebhook(w),
	})
}

func (IncomingWebhook) Type() WebhookType {
	return WebhookTypeIncoming
}

func (w IncomingWebhook) ID() snowflake.Snowflake {
	return w.id
}

func (w IncomingWebhook) Name() string {
	return w.name
}

func (w IncomingWebhook) Avatar() *string {
	return w.avatar
}

func (w IncomingWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if w.Avatar() == nil {
		return w.DefaultAvatarURL(opts...)
	}
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (w IncomingWebhook) AvatarURL(opts ...CDNOpt) *string {
	return formatAssetURL(route.UserAvatar, opts, w.ID(), w.Avatar())
}

func (w IncomingWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	if avatar := formatAssetURL(route.DefaultUserAvatar, opts, 0); avatar != nil {
		return *avatar
	}
	return ""
}

func (w IncomingWebhook) URL() string {
	if compiledRoute, err := route.WebhookURL.Compile(nil, w.ID(), w.Token); err == nil {
		return compiledRoute.URL()
	}
	return ""
}

func (IncomingWebhook) webhook() {}

var _ Webhook = (*ChannelFollowerWebhook)(nil)

type ChannelFollowerWebhook struct {
	id            snowflake.Snowflake
	name          string
	avatar        *string
	ChannelID     snowflake.Snowflake  `json:"channel_id"`
	GuildID       snowflake.Snowflake  `json:"guild_id"`
	SourceGuild   WebhookSourceGuild   `json:"source_guild"`
	SourceChannel WebhookSourceChannel `json:"source_channel"`
	User          User                 `json:"user"`
}

func (w *ChannelFollowerWebhook) UnmarshalJSON(data []byte) error {
	type channelFollowerWebhook ChannelFollowerWebhook
	var v struct {
		ID     snowflake.Snowflake `json:"id"`
		Name   string              `json:"name"`
		Avatar *string             `json:"avatar"`
		channelFollowerWebhook
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = ChannelFollowerWebhook(v.channelFollowerWebhook)
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	return nil
}

func (w ChannelFollowerWebhook) MarshalJSON() ([]byte, error) {
	type channelFollowerWebhook ChannelFollowerWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		channelFollowerWebhook
	}{
		Type:                   w.Type(),
		channelFollowerWebhook: channelFollowerWebhook(w),
	})
}

func (ChannelFollowerWebhook) Type() WebhookType {
	return WebhookTypeChannelFollower
}

func (w ChannelFollowerWebhook) ID() snowflake.Snowflake {
	return w.id
}

func (w ChannelFollowerWebhook) Name() string {
	return w.name
}

func (w ChannelFollowerWebhook) Avatar() *string {
	return w.avatar
}

func (w ChannelFollowerWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if w.Avatar() == nil {
		return w.DefaultAvatarURL(opts...)
	}
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (w ChannelFollowerWebhook) AvatarURL(opts ...CDNOpt) *string {
	return formatAssetURL(route.UserAvatar, opts, w.ID(), w.Avatar())
}

func (w ChannelFollowerWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	if avatar := formatAssetURL(route.DefaultUserAvatar, opts, 0); avatar != nil {
		return *avatar
	}
	return ""
}

func (ChannelFollowerWebhook) webhook() {}

var _ Webhook = (*ApplicationWebhook)(nil)

type ApplicationWebhook struct {
	id            snowflake.Snowflake
	name          string
	avatar        *string
	ApplicationID snowflake.Snowflake
}

func (w *ApplicationWebhook) UnmarshalJSON(data []byte) error {
	var v struct {
		ID            snowflake.Snowflake `json:"id"`
		Name          string              `json:"name"`
		Avatar        *string             `json:"avatar"`
		ApplicationID snowflake.Snowflake `json:"application_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	w.id = v.ID
	w.name = v.Name
	w.avatar = v.Avatar
	w.ApplicationID = v.ApplicationID
	return nil
}

func (w ApplicationWebhook) MarshalJSON() ([]byte, error) {
	type applicationWebhook ApplicationWebhook
	return json.Marshal(struct {
		Type WebhookType `json:"type"`
		applicationWebhook
	}{
		Type:               w.Type(),
		applicationWebhook: applicationWebhook(w),
	})
}

func (ApplicationWebhook) Type() WebhookType {
	return WebhookTypeApplication
}

func (w ApplicationWebhook) ID() snowflake.Snowflake {
	return w.id
}

func (w ApplicationWebhook) Name() string {
	return w.name
}

func (w ApplicationWebhook) Avatar() *string {
	return w.avatar
}

func (w ApplicationWebhook) EffectiveAvatarURL(opts ...CDNOpt) string {
	if w.Avatar() == nil {
		return w.DefaultAvatarURL(opts...)
	}
	if avatar := w.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (w ApplicationWebhook) AvatarURL(opts ...CDNOpt) *string {
	return formatAssetURL(route.UserAvatar, opts, w.ID(), w.Avatar())
}

func (w ApplicationWebhook) DefaultAvatarURL(opts ...CDNOpt) string {
	if avatar := formatAssetURL(route.DefaultUserAvatar, opts, 0); avatar != nil {
		return *avatar
	}
	return ""
}

func (ApplicationWebhook) webhook() {}

type WebhookSourceGuild struct {
	ID   snowflake.Snowflake  `json:"id"`
	Name string               `json:"name"`
	Icon *json.Nullable[Icon] `json:"icon"`
}

type WebhookSourceChannel struct {
	ID   snowflake.Snowflake `json:"id"`
	Name string              `json:"name"`
}

// WebhookCreate is used to create a Webhook
type WebhookCreate struct {
	Name   string `json:"name"`
	Avatar *Icon  `json:"avatar,omitempty"`
}

// WebhookUpdate is used to update a Webhook
type WebhookUpdate struct {
	Name      *string              `json:"name,omitempty"`
	Avatar    *json.Nullable[Icon] `json:"avatar,omitempty"`
	ChannelID *snowflake.Snowflake `json:"channel_id"`
}

// WebhookUpdateWithToken is used to update a Webhook with the token
type WebhookUpdateWithToken struct {
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}
