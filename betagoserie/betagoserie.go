package betagoserie

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	url_net "net/url"
	"strings"
)

// TODO : put const into config file ( ENV VARIABLE ? )
// const apiKey = "ee7422ce11a2"
// const apiSecret = "d2e555a996fb64e49febb5adb9d1c818"
// const login = "d2e555a996fb64e49febb5adb9d1c818"
// const password = "d2e555a996fb64e49febb5adb9d1c818"

// TODO : handle
//  - Episode Watched
//  - Episode UnWatched
//  - Episode List
//  - Episode Show
//  - Show Display
//  - Show Episodes
//  - Show List

// TODO : AddComment

const baseUrl = "https://api.betaseries.com/"

type BetaClient struct {
	ApiKey string
	Token  string
}

// token is a struct return by betaseries
type tokenStruct struct {
	User struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Xp        int    `json:"xp"`
		InAccount bool   `json:"in_account"`
	} `json:"user"`
	Token  string        `json:"token"`
	Hash   string        `json:"hash"`
	Errors []interface{} `json:"errors"`
}

type episodeStruct struct {
	ID        int    `json:"id"`
	TvdbID    int    `json:"thetvdb_id"`
	YoutubeID int    `json:"youtube_id"`
	Title     string `json:"title"`
	Season    string `json:"season"`
	Episode   int    `json:"episode"`
	Show      struct {
		ID               int    `json:"id"`
		TvdbID           int    `json:"thetvdb_id"`
		Title            string `json:"title"`
		InAccount        bool   `json:"in_account"`
		Remaining        int    `json:"remaining"`
		MinutesRemaining int    `json:"minutes_remaining"`
		Progress         int    `json:"progress"`
	} `json:"show"`
	Code        string `json:"code"`
	Global      int    `json:"global"`
	Special     int    `json:"special"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Note        struct {
		Total int     `json:"total"`
		Mean  float32 `json:"mean"`
		User  int     `json:"user"`
	} `json:"note"`
	User struct {
		Seen       bool `json:"seen"`
		Downloaded bool `json:"downloaded"`
	} `json:"user"`
	Comments   string `json:"comments"`
	ResoureURL string `json:"resource_url"`
}

type showsStruct struct {
	ID            int             `json:"id"`
	TvdbID        int             `json:"thetvdb_id"`
	ImdbID        string          `json:"imdb_id"`
	Title         string          `json:"title"`
	Remaining     int             `json:"remaining"`
	EpisodeUnseen []episodeStruct `json:"unseen"`
}

type errorStruct struct {
}

// episode is a struct return by betaseries
type episodeListStruct struct {
	Show  []showsStruct `json:"shows"`
	Error []errorStruct `json:"error"`
}

func NewBetaClient(apiKey, login, password string) *BetaClient {

	finished := make(chan bool)

	var bs = &BetaClient{
		ApiKey: apiKey,
		Token:  " ",
	}

	go bs.getAuthToken(login, password, finished)
	<-finished

	return bs
}

func (bs *BetaClient) GetListEpisode() string {
	var url = baseUrl + "episodes/list"
	var params = map[string]string{
		"token": bs.Token,
	}
	result := bs.makeRequest(url, "GET", params)
	// var list episodeListStruct
	// json.NewDecoder(strings.NewReader(result)).Decode(&list)
	return result
}

func (bs *BetaClient) getAuthToken(login, password string, finished chan bool) {
	hasher := md5.New()
	hasher.Write([]byte(password))

	var url = baseUrl + "members/auth"
	var params = map[string]string{
		"login":    login,
		"password": hex.EncodeToString(hasher.Sum(nil)),
	}

	var token tokenStruct
	result := bs.makeRequest(url, "POST", params)
	json.NewDecoder(strings.NewReader(result)).Decode(&token)
	bs.Token = token.Token
	finished <- true
}

// Make the request and return the JSON Data
func (bs *BetaClient) makeRequest(url, urlType string, params map[string]string) string {

	data := url_net.Values{}
	data.Set("client_id", bs.ApiKey)

	for paramKey, paramValue := range params {
		data.Add(paramKey, paramValue)
	}

	// Build the request
	req, err := http.NewRequest(urlType, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("erreur ReadAll: ", err)
	}

	return string(body)

}
