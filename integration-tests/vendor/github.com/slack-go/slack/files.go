package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	// Add here the defaults in the siten
	DEFAULT_FILES_USER        = ""
	DEFAULT_FILES_CHANNEL     = ""
	DEFAULT_FILES_TS_FROM     = 0
	DEFAULT_FILES_TS_TO       = -1
	DEFAULT_FILES_TYPES       = "all"
	DEFAULT_FILES_COUNT       = 100
	DEFAULT_FILES_PAGE        = 1
	DEFAULT_FILES_SHOW_HIDDEN = false
)

// File contains all the information for a file
type File struct {
	ID        string   `json:"id"`
	Created   JSONTime `json:"created"`
	Timestamp JSONTime `json:"timestamp"`

	Name              string `json:"name"`
	Title             string `json:"title"`
	Mimetype          string `json:"mimetype"`
	ImageExifRotation int    `json:"image_exif_rotation"`
	Filetype          string `json:"filetype"`
	PrettyType        string `json:"pretty_type"`
	User              string `json:"user"`

	Mode         string `json:"mode"`
	Editable     bool   `json:"editable"`
	IsExternal   bool   `json:"is_external"`
	ExternalType string `json:"external_type"`

	Size int `json:"size"`

	URL                string `json:"url"`          // Deprecated - never set
	URLDownload        string `json:"url_download"` // Deprecated - never set
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`

	OriginalH   int    `json:"original_h"`
	OriginalW   int    `json:"original_w"`
	Thumb64     string `json:"thumb_64"`
	Thumb80     string `json:"thumb_80"`
	Thumb160    string `json:"thumb_160"`
	Thumb360    string `json:"thumb_360"`
	Thumb360Gif string `json:"thumb_360_gif"`
	Thumb360W   int    `json:"thumb_360_w"`
	Thumb360H   int    `json:"thumb_360_h"`
	Thumb480    string `json:"thumb_480"`
	Thumb480W   int    `json:"thumb_480_w"`
	Thumb480H   int    `json:"thumb_480_h"`
	Thumb720    string `json:"thumb_720"`
	Thumb720W   int    `json:"thumb_720_w"`
	Thumb720H   int    `json:"thumb_720_h"`
	Thumb960    string `json:"thumb_960"`
	Thumb960W   int    `json:"thumb_960_w"`
	Thumb960H   int    `json:"thumb_960_h"`
	Thumb1024   string `json:"thumb_1024"`
	Thumb1024W  int    `json:"thumb_1024_w"`
	Thumb1024H  int    `json:"thumb_1024_h"`

	Permalink       string `json:"permalink"`
	PermalinkPublic string `json:"permalink_public"`

	EditLink         string `json:"edit_link"`
	Preview          string `json:"preview"`
	PreviewHighlight string `json:"preview_highlight"`
	Lines            int    `json:"lines"`
	LinesMore        int    `json:"lines_more"`

	IsPublic        bool     `json:"is_public"`
	PublicURLShared bool     `json:"public_url_shared"`
	Channels        []string `json:"channels"`
	Groups          []string `json:"groups"`
	IMs             []string `json:"ims"`
	InitialComment  Comment  `json:"initial_comment"`
	CommentsCount   int      `json:"comments_count"`
	NumStars        int      `json:"num_stars"`
	IsStarred       bool     `json:"is_starred"`
	Shares          Share    `json:"shares"`
}

type Share struct {
	Public  map[string][]ShareFileInfo `json:"public"`
	Private map[string][]ShareFileInfo `json:"private"`
}

type ShareFileInfo struct {
	ReplyUsers      []string `json:"reply_users"`
	ReplyUsersCount int      `json:"reply_users_count"`
	ReplyCount      int      `json:"reply_count"`
	Ts              string   `json:"ts"`
	ThreadTs        string   `json:"thread_ts"`
	LatestReply     string   `json:"latest_reply"`
	ChannelName     string   `json:"channel_name"`
	TeamID          string   `json:"team_id"`
}

// FileUploadParameters contains all the parameters necessary (including the optional ones) for an UploadFile() request.
//
// There are three ways to upload a file. You can either set Content if file is small, set Reader if file is large,
// or provide a local file path in File to upload it from your filesystem.
//
// Note that when using the Reader option, you *must* specify the Filename, otherwise the Slack API isn't happy.
type FileUploadParameters struct {
	File            string
	Content         string
	Reader          io.Reader
	Filetype        string
	Filename        string
	Title           string
	InitialComment  string
	Channels        []string
	ThreadTimestamp string
}

// GetFilesParameters contains all the parameters necessary (including the optional ones) for a GetFiles() request
type GetFilesParameters struct {
	User          string
	Channel       string
	TimestampFrom JSONTime
	TimestampTo   JSONTime
	Types         string
	Count         int
	Page          int
	ShowHidden    bool
}

// ListFilesParameters contains all the parameters necessary (including the optional ones) for a ListFiles() request
type ListFilesParameters struct {
	Limit   int
	User    string
	Channel string
	Types   string
	Cursor  string
}

type UploadFileV2Parameters struct {
	File            string
	FileSize        int
	Content         string
	Reader          io.Reader
	Filename        string
	Title           string
	InitialComment  string
	Channel         string
	ThreadTimestamp string
	AltTxt          string
	SnippetText     string
}

type getUploadURLExternalParameters struct {
	altText     string
	fileSize    int
	fileName    string
	snippetText string
}

type getUploadURLExternalResponse struct {
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
	SlackResponse
}

type uploadToURLParameters struct {
	UploadURL string
	Reader    io.Reader
	File      string
	Content   string
	Filename  string
}

type FileSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type completeUploadExternalParameters struct {
	title           string
	channel         string
	initialComment  string
	threadTimestamp string
}

type completeUploadExternalResponse struct {
	SlackResponse
	Files []FileSummary `json:"files"`
}

type fileResponseFull struct {
	File     `json:"file"`
	Paging   `json:"paging"`
	Comments []Comment        `json:"comments"`
	Files    []File           `json:"files"`
	Metadata ResponseMetadata `json:"response_metadata"`

	SlackResponse
}

// NewGetFilesParameters provides an instance of GetFilesParameters with all the sane default values set
func NewGetFilesParameters() GetFilesParameters {
	return GetFilesParameters{
		User:          DEFAULT_FILES_USER,
		Channel:       DEFAULT_FILES_CHANNEL,
		TimestampFrom: DEFAULT_FILES_TS_FROM,
		TimestampTo:   DEFAULT_FILES_TS_TO,
		Types:         DEFAULT_FILES_TYPES,
		Count:         DEFAULT_FILES_COUNT,
		Page:          DEFAULT_FILES_PAGE,
		ShowHidden:    DEFAULT_FILES_SHOW_HIDDEN,
	}
}

func (api *Client) fileRequest(ctx context.Context, path string, values url.Values) (*fileResponseFull, error) {
	response := &fileResponseFull{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// GetFileInfo retrieves a file and related comments
func (api *Client) GetFileInfo(fileID string, count, page int) (*File, []Comment, *Paging, error) {
	return api.GetFileInfoContext(context.Background(), fileID, count, page)
}

// GetFileInfoContext retrieves a file and related comments with a custom context
func (api *Client) GetFileInfoContext(ctx context.Context, fileID string, count, page int) (*File, []Comment, *Paging, error) {
	values := url.Values{
		"token": {api.token},
		"file":  {fileID},
		"count": {strconv.Itoa(count)},
		"page":  {strconv.Itoa(page)},
	}

	response, err := api.fileRequest(ctx, "files.info", values)
	if err != nil {
		return nil, nil, nil, err
	}
	return &response.File, response.Comments, &response.Paging, nil
}

// GetFile retreives a given file from its private download URL
func (api *Client) GetFile(downloadURL string, writer io.Writer) error {
	return api.GetFileContext(context.Background(), downloadURL, writer)
}

// GetFileContext retreives a given file from its private download URL with a custom context
//
// For more details, see GetFile documentation.
func (api *Client) GetFileContext(ctx context.Context, downloadURL string, writer io.Writer) error {
	return downloadFile(ctx, api.httpclient, api.token, downloadURL, writer, api)
}

// GetFiles retrieves all files according to the parameters given
func (api *Client) GetFiles(params GetFilesParameters) ([]File, *Paging, error) {
	return api.GetFilesContext(context.Background(), params)
}

// GetFilesContext retrieves all files according to the parameters given with a custom context
func (api *Client) GetFilesContext(ctx context.Context, params GetFilesParameters) ([]File, *Paging, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.User != DEFAULT_FILES_USER {
		values.Add("user", params.User)
	}
	if params.Channel != DEFAULT_FILES_CHANNEL {
		values.Add("channel", params.Channel)
	}
	if params.TimestampFrom != DEFAULT_FILES_TS_FROM {
		values.Add("ts_from", strconv.FormatInt(int64(params.TimestampFrom), 10))
	}
	if params.TimestampTo != DEFAULT_FILES_TS_TO {
		values.Add("ts_to", strconv.FormatInt(int64(params.TimestampTo), 10))
	}
	if params.Types != DEFAULT_FILES_TYPES {
		values.Add("types", params.Types)
	}
	if params.Count != DEFAULT_FILES_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_FILES_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	if params.ShowHidden != DEFAULT_FILES_SHOW_HIDDEN {
		values.Add("show_files_hidden_by_limit", strconv.FormatBool(params.ShowHidden))
	}

	response, err := api.fileRequest(ctx, "files.list", values)
	if err != nil {
		return nil, nil, err
	}
	return response.Files, &response.Paging, nil
}

// ListFiles retrieves all files according to the parameters given. Uses cursor based pagination.
func (api *Client) ListFiles(params ListFilesParameters) ([]File, *ListFilesParameters, error) {
	return api.ListFilesContext(context.Background(), params)
}

// ListFilesContext retrieves all files according to the parameters given with a custom context.
//
// For more details, see ListFiles documentation.
func (api *Client) ListFilesContext(ctx context.Context, params ListFilesParameters) ([]File, *ListFilesParameters, error) {
	values := url.Values{
		"token": {api.token},
	}

	if params.User != DEFAULT_FILES_USER {
		values.Add("user", params.User)
	}
	if params.Channel != DEFAULT_FILES_CHANNEL {
		values.Add("channel", params.Channel)
	}
	if params.Limit != DEFAULT_FILES_COUNT {
		values.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		values.Add("cursor", params.Cursor)
	}

	response, err := api.fileRequest(ctx, "files.list", values)
	if err != nil {
		return nil, nil, err
	}

	params.Cursor = response.Metadata.Cursor

	return response.Files, &params, nil
}

// UploadFile uploads a file
func (api *Client) UploadFile(params FileUploadParameters) (file *File, err error) {
	return api.UploadFileContext(context.Background(), params)
}

// UploadFileContext uploads a file and setting a custom context
func (api *Client) UploadFileContext(ctx context.Context, params FileUploadParameters) (file *File, err error) {
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTestContext(ctx)
	if err != nil {
		return nil, err
	}
	response := &fileResponseFull{}
	values := url.Values{}
	if params.Filetype != "" {
		values.Add("filetype", params.Filetype)
	}
	if params.Filename != "" {
		values.Add("filename", params.Filename)
	}
	if params.Title != "" {
		values.Add("title", params.Title)
	}
	if params.InitialComment != "" {
		values.Add("initial_comment", params.InitialComment)
	}
	if params.ThreadTimestamp != "" {
		values.Add("thread_ts", params.ThreadTimestamp)
	}
	if len(params.Channels) != 0 {
		values.Add("channels", strings.Join(params.Channels, ","))
	}
	if params.Content != "" {
		values.Add("content", params.Content)
		values.Add("token", api.token)
		err = api.postMethod(ctx, "files.upload", values, response)
	} else if params.File != "" {
		err = postLocalWithMultipartResponse(ctx, api.httpclient, api.endpoint+"files.upload", params.File, "file", api.token, values, response, api)
	} else if params.Reader != nil {
		if params.Filename == "" {
			return nil, fmt.Errorf("files.upload: FileUploadParameters.Filename is mandatory when using FileUploadParameters.Reader")
		}
		err = postWithMultipartResponse(ctx, api.httpclient, api.endpoint+"files.upload", params.Filename, "file", api.token, values, params.Reader, response, api)
	}

	if err != nil {
		return nil, err
	}

	return &response.File, response.Err()
}

// DeleteFileComment deletes a file's comment
func (api *Client) DeleteFileComment(commentID, fileID string) error {
	return api.DeleteFileCommentContext(context.Background(), fileID, commentID)
}

// DeleteFileCommentContext deletes a file's comment with a custom context
func (api *Client) DeleteFileCommentContext(ctx context.Context, fileID, commentID string) (err error) {
	if fileID == "" || commentID == "" {
		return ErrParametersMissing
	}

	values := url.Values{
		"token": {api.token},
		"file":  {fileID},
		"id":    {commentID},
	}
	_, err = api.fileRequest(ctx, "files.comments.delete", values)
	return err
}

// DeleteFile deletes a file
func (api *Client) DeleteFile(fileID string) error {
	return api.DeleteFileContext(context.Background(), fileID)
}

// DeleteFileContext deletes a file with a custom context
func (api *Client) DeleteFileContext(ctx context.Context, fileID string) (err error) {
	values := url.Values{
		"token": {api.token},
		"file":  {fileID},
	}

	_, err = api.fileRequest(ctx, "files.delete", values)
	return err
}

// RevokeFilePublicURL disables public/external sharing for a file
func (api *Client) RevokeFilePublicURL(fileID string) (*File, error) {
	return api.RevokeFilePublicURLContext(context.Background(), fileID)
}

// RevokeFilePublicURLContext disables public/external sharing for a file with a custom context
func (api *Client) RevokeFilePublicURLContext(ctx context.Context, fileID string) (*File, error) {
	values := url.Values{
		"token": {api.token},
		"file":  {fileID},
	}

	response, err := api.fileRequest(ctx, "files.revokePublicURL", values)
	if err != nil {
		return nil, err
	}
	return &response.File, nil
}

// ShareFilePublicURL enabled public/external sharing for a file
func (api *Client) ShareFilePublicURL(fileID string) (*File, []Comment, *Paging, error) {
	return api.ShareFilePublicURLContext(context.Background(), fileID)
}

// ShareFilePublicURLContext enabled public/external sharing for a file with a custom context
func (api *Client) ShareFilePublicURLContext(ctx context.Context, fileID string) (*File, []Comment, *Paging, error) {
	values := url.Values{
		"token": {api.token},
		"file":  {fileID},
	}

	response, err := api.fileRequest(ctx, "files.sharedPublicURL", values)
	if err != nil {
		return nil, nil, nil, err
	}
	return &response.File, response.Comments, &response.Paging, nil
}

// getUploadURLExternal gets a URL and fileID from slack which can later be used to upload a file
func (api *Client) getUploadURLExternal(ctx context.Context, params getUploadURLExternalParameters) (*getUploadURLExternalResponse, error) {
	values := url.Values{
		"token":    {api.token},
		"filename": {params.fileName},
		"length":   {strconv.Itoa(params.fileSize)},
	}
	if params.altText != "" {
		values.Add("initial_comment", params.altText)
	}
	if params.snippetText != "" {
		values.Add("thread_ts", params.snippetText)
	}
	response := &getUploadURLExternalResponse{}
	err := api.postMethod(ctx, "files.getUploadURLExternal", values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

// uploadToURL uploads the file to the provided URL using post method
func (api *Client) uploadToURL(ctx context.Context, params uploadToURLParameters) (err error) {
	values := url.Values{}
	if params.Content != "" {
		values.Add("content", params.Content)
		values.Add("token", api.token)
		err = postForm(ctx, api.httpclient, params.UploadURL, values, nil, api)
	} else if params.File != "" {
		err = postLocalWithMultipartResponse(ctx, api.httpclient, params.UploadURL, params.File, "file", api.token, values, nil, api)
	} else if params.Reader != nil {
		err = postWithMultipartResponse(ctx, api.httpclient, params.UploadURL, params.Filename, "file", api.token, values, params.Reader, nil, api)
	}
	return err
}

// completeUploadExternal once files are uploaded, this completes the upload and shares it to the specified channel
func (api *Client) completeUploadExternal(ctx context.Context, fileID string, params completeUploadExternalParameters) (file *completeUploadExternalResponse, err error) {
	request := []FileSummary{{ID: fileID, Title: params.title}}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"token":      {api.token},
		"files":      {string(requestBytes)},
		"channel_id": {params.channel},
	}

	if params.initialComment != "" {
		values.Add("initial_comment", params.initialComment)
	}
	if params.threadTimestamp != "" {
		values.Add("thread_ts", params.threadTimestamp)
	}
	response := &completeUploadExternalResponse{}
	err = api.postMethod(ctx, "files.completeUploadExternal", values, response)
	if err != nil {
		return nil, err
	}
	if response.Err() != nil {
		return nil, response.Err()
	}
	return response, nil
}

// UploadFileV2 uploads file to a given slack channel using 3 steps -
//  1. Get an upload URL using files.getUploadURLExternal API
//  2. Send the file as a post to the URL provided by slack
//  3. Complete the upload and share it to the specified channel using files.completeUploadExternal
func (api *Client) UploadFileV2(params UploadFileV2Parameters) (*FileSummary, error) {
	return api.UploadFileV2Context(context.Background(), params)
}

// UploadFileV2 uploads file to a given slack channel using 3 steps with a custom context -
//  1. Get an upload URL using files.getUploadURLExternal API
//  2. Send the file as a post to the URL provided by slack
//  3. Complete the upload and share it to the specified channel using files.completeUploadExternal
func (api *Client) UploadFileV2Context(ctx context.Context, params UploadFileV2Parameters) (file *FileSummary, err error) {
	if params.Filename == "" {
		return nil, fmt.Errorf("file.upload.v2: filename cannot be empty")
	}
	if params.FileSize == 0 {
		return nil, fmt.Errorf("file.upload.v2: file size cannot be 0")
	}
	if params.Channel == "" {
		return nil, fmt.Errorf("file.upload.v2: channel cannot be empty")
	}
	u, err := api.getUploadURLExternal(ctx, getUploadURLExternalParameters{
		altText:     params.AltTxt,
		fileName:    params.Filename,
		fileSize:    params.FileSize,
		snippetText: params.SnippetText,
	})
	if err != nil {
		return nil, err
	}

	err = api.uploadToURL(ctx, uploadToURLParameters{
		UploadURL: u.UploadURL,
		Reader:    params.Reader,
		File:      params.File,
		Content:   params.Content,
		Filename:  params.Filename,
	})
	if err != nil {
		return nil, err
	}

	c, err := api.completeUploadExternal(ctx, u.FileID, completeUploadExternalParameters{
		title:           params.Title,
		channel:         params.Channel,
		initialComment:  params.InitialComment,
		threadTimestamp: params.ThreadTimestamp,
	})
	if err != nil {
		return nil, err
	}
	if len(c.Files) != 1 {
		return nil, fmt.Errorf("file.upload.v2: something went wrong; received %d files instead of 1", len(c.Files))
	}

	return &c.Files[0], nil
}
