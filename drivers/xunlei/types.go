package xunlei

import (
	"fmt"
	"time"
)

type Erron struct {
	ErrorCode        int64  `json:"error_code"`
	ErrorMsg         string `json:"error"`
	ErrorDescription string `json:"error_description"`
	//	ErrorDetails   interface{} `json:"error_details"`
}

func (e *Erron) HasError() bool {
	return e.ErrorCode != 0 || e.ErrorMsg != "" || e.ErrorDescription != ""
}

func (e *Erron) Error() string {
	return fmt.Sprintf("ErrorCode: %d ,Error: %s ,ErrorDescription: %s ", e.ErrorCode, e.ErrorMsg, e.ErrorDescription)
}

/*
* 验证码Token
**/
type CaptchaTokenRequest struct {
	Action       string            `json:"action"`
	CaptchaToken string            `json:"captcha_token"`
	ClientID     string            `json:"client_id"`
	DeviceID     string            `json:"device_id"`
	Meta         map[string]string `json:"meta"`
	RedirectUri  string            `json:"redirect_uri"`
}

type CaptchaTokenResponse struct {
	CaptchaToken string `json:"captcha_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Url          string `json:"url"`
}

/*
* 登录
**/
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`

	Sub    string `json:"sub"`
	UserID string `json:"user_id"`
}

func (t *TokenResponse) Token() string {
	return fmt.Sprint(t.TokenType, " ", t.AccessToken)
}

type SignInRequest struct {
	CaptchaToken string `json:"captcha_token"`

	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	Username string `json:"username"`
	Password string `json:"password"`
}

/*
* 文件
**/
type FileList struct {
	Kind            string  `json:"kind"`
	NextPageToken   string  `json:"next_page_token"`
	Files           []Files `json:"files"`
	Version         string  `json:"version"`
	VersionOutdated bool    `json:"version_outdated"`
}

type Files struct {
	Kind           string     `json:"kind"`
	ID             string     `json:"id"`
	ParentID       string     `json:"parent_id"`
	Name           string     `json:"name"`
	UserID         string     `json:"user_id"`
	Size           string     `json:"size"`
	Revision       string     `json:"revision"`
	FileExtension  string     `json:"file_extension"`
	MimeType       string     `json:"mime_type"`
	Starred        bool       `json:"starred"`
	WebContentLink string     `json:"web_content_link"`
	CreatedTime    *time.Time `json:"created_time"`
	ModifiedTime   *time.Time `json:"modified_time"`
	IconLink       string     `json:"icon_link"`
	ThumbnailLink  string     `json:"thumbnail_link"`
	Md5Checksum    string     `json:"md5_checksum"`
	Hash           string     `json:"hash"`
	//Links          struct{}   `json:"links"`
	Phase string `json:"phase"`
	Audit struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Title   string `json:"title"`
	} `json:"audit"`
	/* Medias []struct {
		Category       string      `json:"category"`
		IconLink       string      `json:"icon_link"`
		IsDefault      bool        `json:"is_default"`
		IsOrigin       bool        `json:"is_origin"`
		IsVisible      bool        `json:"is_visible"`
		//Link           interface{} `json:"link"`
		MediaID        string      `json:"media_id"`
		MediaName      string      `json:"media_name"`
		NeedMoreQuota  bool        `json:"need_more_quota"`
		Priority       int         `json:"priority"`
		RedirectLink   string      `json:"redirect_link"`
		ResolutionName string      `json:"resolution_name"`
		Video          struct {
			AudioCodec string `json:"audio_codec"`
			BitRate    int    `json:"bit_rate"`
			Duration   int    `json:"duration"`
			FrameRate  int    `json:"frame_rate"`
			Height     int    `json:"height"`
			VideoCodec string `json:"video_codec"`
			VideoType  string `json:"video_type"`
			Width      int    `json:"width"`
		} `json:"video"`
		VipTypes []string `json:"vip_types"`
	} `json:"medias"` */
	Trashed     bool   `json:"trashed"`
	DeleteTime  string `json:"delete_time"`
	OriginalURL string `json:"original_url"`
	//Params            struct{} `json:"params"`
	OriginalFileIndex int    `json:"original_file_index"`
	Space             string `json:"space"`
	//Apps              []interface{} `json:"apps"`
	Writable   bool   `json:"writable"`
	FolderType string `json:"folder_type"`
	//Collection interface{} `json:"collection"`
}

/*
* 上传
**/
type UploadTaskResponse struct {
	UploadType string `json:"upload_type"`

	/*//UPLOAD_TYPE_FORM
	Form struct {
		//Headers struct{} `json:"headers"`
		Kind       string `json:"kind"`
		Method     string `json:"method"`
		MultiParts struct {
			OSSAccessKeyID string `json:"OSSAccessKeyId"`
			Signature      string `json:"Signature"`
			Callback       string `json:"callback"`
			Key            string `json:"key"`
			Policy         string `json:"policy"`
			XUserData      string `json:"x:user_data"`
		} `json:"multi_parts"`
		URL string `json:"url"`
	} `json:"form"`*/

	//UPLOAD_TYPE_RESUMABLE
	Resumable struct {
		Kind   string `json:"kind"`
		Params struct {
			AccessKeyID     string    `json:"access_key_id"`
			AccessKeySecret string    `json:"access_key_secret"`
			Bucket          string    `json:"bucket"`
			Endpoint        string    `json:"endpoint"`
			Expiration      time.Time `json:"expiration"`
			Key             string    `json:"key"`
			SecurityToken   string    `json:"security_token"`
		} `json:"params"`
		Provider string `json:"provider"`
	} `json:"resumable"`

	File Files `json:"file"`
}
