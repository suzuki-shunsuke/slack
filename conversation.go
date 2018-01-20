package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// Conversation is the foundation for IM and BaseGroupConversation
type conversation struct {
	ID                 string   `json:"id"`
	Created            JSONTime `json:"created"`
	IsOpen             bool     `json:"is_open"`
	LastRead           string   `json:"last_read,omitempty"`
	Latest             *Message `json:"latest,omitempty"`
	UnreadCount        int      `json:"unread_count,omitempty"`
	UnreadCountDisplay int      `json:"unread_count_display,omitempty"`
}

// GroupConversation is the foundation for Group and Channel
type groupConversation struct {
	conversation
	Name       string   `json:"name"`
	Creator    string   `json:"creator"`
	IsArchived bool     `json:"is_archived"`
	Members    []string `json:"members"`
	Topic      Topic    `json:"topic"`
	Purpose    Purpose  `json:"purpose"`
}

// Topic contains information about the topic
type Topic struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

// Purpose contains information about the purpose
type Purpose struct {
	Value   string   `json:"value"`
	Creator string   `json:"creator"`
	LastSet JSONTime `json:"last_set"`
}

type GetUsersInConversationParameters struct {
	ChannelID string
	Cursor    string
	Limit     int
}

type responseMetaData struct {
	NextCursor string `json:"next_cursor"`
}

// GetUsersInConversation returns the list of users in a conversation
func (api *Client) GetUsersInConversation(params *GetUsersInConversationParameters) ([]string, string, error) {
	return api.GetUsersInConversationContext(context.Background(), params)
}

// GetUsersInConversation returns the list of users in a conversation with a custom context
func (api *Client) GetUsersInConversationContext(ctx context.Context, params *GetUsersInConversationParameters) ([]string, string, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {params.ChannelID},
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}
	if params.Limit != 0 {
		values.Add("limit", string(params.Limit))
	}
	response := struct {
		Members          []string         `json:"members"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
		SlackResponse
	}{}
	err := post(ctx, api.httpclient, "conversations.members", values, &response, api.debug)
	if err != nil {
		return nil, "", err
	}
	if !response.Ok {
		return nil, "", errors.New(response.Error)
	}
	return response.Members, response.ResponseMetaData.NextCursor, nil
}

// ArchiveConversation archives a conversation
func (api *Client) ArchiveConversation(channelID string) error {
	return api.ArchiveConversationContext(context.Background(), channelID)
}

// ArchiveConversationContext archives a conversation with a custom context
func (api *Client) ArchiveConversationContext(ctx context.Context, channelID string) error {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
	}
	response := SlackResponse{}
	err := post(ctx, api.httpclient, "conversations.archive", values, &response, api.debug)
	if err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// UnArchiveConversation reverses conversation archival
func (api *Client) UnArchiveConversation(channelID string) error {
	return api.UnArchiveConversationContext(context.Background(), channelID)
}

// UnArchiveConversationContext reverses conversation archival with a custom context
func (api *Client) UnArchiveConversationContext(ctx context.Context, channelID string) error {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
	}
	response := SlackResponse{}
	err := post(ctx, api.httpclient, "conversations.unarchive", values, &response, api.debug)
	if err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// SetTopicOfConversation sets the topic for a conversation
func (api *Client) SetTopicOfConversation(channelID, topic string) (*Channel, error) {
	return api.SetTopicOfConversationContext(context.Background(), channelID, topic)
}

// SetTopicOfConversationContext sets the topic for a conversation with a custom context
func (api *Client) SetTopicOfConversationContext(ctx context.Context, channelID, topic string) (*Channel, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
		"topic":   {topic},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, api.httpclient, "conversations.setTopic", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// SetPurposeOfConversation sets the purpose for a conversation
func (api *Client) SetPurposeOfConversation(channelID, purpose string) (*Channel, error) {
	return api.SetPurposeOfConversationContext(context.Background(), channelID, purpose)
}

// SetPurposeOfConversationContext sets the purpose for a conversation with a custom context
func (api *Client) SetPurposeOfConversationContext(ctx context.Context, channelID, purpose string) (*Channel, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
		"purpose": {purpose},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, api.httpclient, "conversations.setPurpose", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// RenameConversation renames a conversation
func (api *Client) RenameConversation(channelID, channelName string) (*Channel, error) {
	return api.RenameConversationContext(context.Background(), channelID, channelName)
}

// RenameConversationContext renames a conversation with a custom context
func (api *Client) RenameConversationContext(ctx context.Context, channelID, channelName string) (*Channel, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
		"name":    {channelName},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, api.httpclient, "conversations.rename", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// InviteUsersToConversation invites users to a channel
func (api *Client) InviteUsersToConversation(channelID string, users []string) (*Channel, error) {
	return api.InviteUsersToConversationContext(context.Background(), channelID, users)
}

// InviteUsersToConversationContext invites users to a channel with a custom context
func (api *Client) InviteUsersToConversationContext(ctx context.Context, channelID string, users []string) (*Channel, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
		"users":   {strings.Join(users, ",")},
	}
	response := struct {
		SlackResponse
		Channel *Channel `json:"channel"`
	}{}
	err := post(ctx, api.httpclient, "conversations.invite", values, &response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Channel, nil
}

// KickUserFromConversation removes a user from a conversation
func (api *Client) KickUserFromConversation(channelID string, user string) error {
	return api.KickUserFromConversationContext(context.Background(), channelID, user)
}

// KickUserFromConversationContext removes a user from a conversation with a custom context
func (api *Client) KickUserFromConversationContext(ctx context.Context, channelID string, user string) error {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
		"user":    {user},
	}
	response := SlackResponse{}
	err := post(ctx, api.httpclient, "conversations.kick", values, &response, api.debug)
	if err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// CloseConversation closes a direct message or multi-person direct message
func (api *Client) CloseConversation(channelID string) (bool, bool, error) {
	return api.CloseConversationContext(context.Background(), channelID)
}

// CloseConversationContext closes a direct message or multi-person direct message with a custom context
func (api *Client) CloseConversationContext(ctx context.Context, channelID string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.token},
		"channel": {channelID},
	}
	response := struct {
		SlackResponse
		NoOp          bool `json:"no_op"`
		AlreadyClosed bool `json:"already_closed"`
	}{}

	err := post(ctx, api.httpclient, "conversations.close", values, &response, api.debug)
	if err != nil {
		return false, false, err
	}
	if !response.Ok {
		return false, false, errors.New(response.Error)
	}
	return response.NoOp, response.AlreadyClosed, nil
}

// CreateConversation initiates a public or private channel-based conversation
func (api *Client) CreateConversation(channelName string, isPrivate bool) (*Channel, error) {
	return api.CreateConversationContext(context.Background(), channelName, isPrivate)
}

// CreateConversationContext initiates a public or private channel-based conversation with a custom context
func (api *Client) CreateConversationContext(ctx context.Context, channelName string, isPrivate bool) (*Channel, error) {
	values := url.Values{
		"token":      {api.token},
		"name":       {channelName},
		"is_private": {strconv.FormatBool(isPrivate)},
	}
	response, err := channelRequest(
		ctx, api.httpclient, "conversations.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return &response.Channel, nil
}

// GetConversationInfo retrieves information about a conversation
func (api *Client) GetConversationInfo(channelID string, includeLocale bool) (*Channel, error) {
	return api.GetConversationInfoContext(context.Background(), channelID, includeLocale)
}

// GetConversationInfoContext retrieves information about a conversation with a custom context
func (api *Client) GetConversationInfoContext(ctx context.Context, channelID string, includeLocale bool) (*Channel, error) {
	values := url.Values{
		"token":          {api.token},
		"channel":        {channelID},
		"include_locale": {strconv.FormatBool(includeLocale)},
	}
	response, err := channelRequest(
		ctx, api.httpclient, "conversations.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return &response.Channel, nil
}
