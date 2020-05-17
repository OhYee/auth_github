package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/OhYee/rainbow/errors"
)

// Connect Github OAuth APP
type Connect struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

// New 新建一个 QQ互联应用，需要传入已申请成功的应用参数
func New(clientID, clientSecret, redirectURI string) *Connect {
	return &Connect{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}

}

/*
LoginPage 返回跳转的登录页面，需要传入一个随机的 state 用于验证身份

https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/
*/
func (conn *Connect) LoginPage(state string, signedUser bool, scope ...string) string {
	params := make(url.Values)
	params.Add("client_id", conn.clientID)
	params.Add("redirect_uri", conn.redirectURI)
	params.Add("scope", strings.Join(scope, " "))
	params.Add("state", state)
	if signedUser {
		params.Add("allow_signup", "true")
	} else {
		params.Add("allow_signup", "false")
	}
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?%s",
		params.Encode(),
	)
}

/*
Auth 根据登录回调页面返回的 code 获取用户 token

https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
*/
func (conn *Connect) Auth(code, state string) (token string, err error) {
	params := make(url.Values)
	params.Add("client_id", conn.clientID)
	params.Add("client_secret", conn.clientSecret)
	params.Add("code", code)
	params.Add("redirect_uri", conn.redirectURI)
	params.Add("state", state)

	resp, err := http.Post(
		fmt.Sprintf("https://github.com/login/oauth/access_token?%s", params.Encode()),
		"application/json",
		&bytes.Buffer{},
	)
	if err != nil {
		errors.Wrapper(&err)
		return
	}

	var b []byte
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		errors.Wrapper(&err)
		return
	}

	urls, err := url.ParseQuery(string(b))
	if err != nil {
		errors.Wrapper(&err)
		return
	}
	token = urls.Get("access_token")
	return
}

// UserInfo get_user_info api response
type UserInfo struct {
	Message           string `json:"message"`
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	Avatar            string `json:"avatar_url"`
	Gravatar          string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Company           string `json:"company"`
	Blog              string `json:"blog"`
	Location          string `json:"location"`
	Email             string `json:"email"`
	Hireable          string `json:"hireable"`
	Bio               string `json:"bio"`
	PublicRepos       int64  `json:"public_repos"`
	PublicGists       int64  `json:"public_gists"`
	Followers         int64  `json:"followers"`
	Following         int64  `json:"following"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

/*
Info 获取用户信息

https://wiki.connect.qq.com/get_user_info

获取登录用户在QQ空间的信息，包括昵称、头像、性别及黄钻信息（包括黄钻等级、是否年费黄钻等）。
*/
func (conn *Connect) Info(token string) (res UserInfo, err error) {
	params := make(url.Values)
	params.Add("access_token", token)

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return
	}
	// Using Authorization Header instead of query param
	// https://developer.github.com/changes/2020-02-10-deprecating-auth-through-query-param/
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		return
	}

	var b []byte
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		errors.Wrapper(&err)
		return
	}

	if err = json.Unmarshal(b, &res); err != nil {
		errors.Wrapper(&err)
		return
	}

	if res.Message != "" {
		err = errors.New(res.Message)
		return
	}

	return
}
