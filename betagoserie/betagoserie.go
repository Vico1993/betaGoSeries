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
	"strconv"
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
//  - Episode Scrapper
//  - Episode Search

// TODO : AddComment

const baseURL = "https://api.betaseries.com/"

// BetaClient Default struct with 2 indispensable parameter
type BetaClient struct {
	APIKey string
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

// type episodeStruct struct {
// 	ID        int    `json:"id"`
// 	TvdbID    int    `json:"thetvdb_id"`
// 	YoutubeID int    `json:"youtube_id"`
// 	Title     string `json:"title"`
// 	Season    string `json:"season"`
// 	Episode   int    `json:"episode"`
// 	Show      struct {
// 		ID               int    `json:"id"`
// 		TvdbID           int    `json:"thetvdb_id"`
// 		Title            string `json:"title"`
// 		InAccount        bool   `json:"in_account"`
// 		Remaining        int    `json:"remaining"`
// 		MinutesRemaining int    `json:"minutes_remaining"`
// 		Progress         int    `json:"progress"`
// 	} `json:"show"`
// 	Code        string `json:"code"`
// 	Global      int    `json:"global"`
// 	Special     int    `json:"special"`
// 	Description string `json:"description"`
// 	Date        string `json:"date"`
// 	Note        struct {
// 		Total int     `json:"total"`
// 		Mean  float32 `json:"mean"`
// 		User  int     `json:"user"`
// 	} `json:"note"`
// 	User struct {
// 		Seen       bool `json:"seen"`
// 		Downloaded bool `json:"downloaded"`
// 	} `json:"user"`
// 	Comments   string `json:"comments"`
// 	ResoureURL string `json:"resource_url"`
// }

// type showsStruct struct {
// 	ID            int             `json:"id"`
// 	TvdbID        int             `json:"thetvdb_id"`
// 	ImdbID        string          `json:"imdb_id"`
// 	Title         string          `json:"title"`
// 	Remaining     int             `json:"remaining"`
// 	EpisodeUnseen []episodeStruct `json:"unseen"`
// }

type errorStruct struct {
}

// episode is a struct return by betaseries
// type episodeListStruct struct {
// 	Show  []showsStruct `json:"shows"`
// 	Error []errorStruct `json:"error"`
// }

// NewBetaClient create or new client to interact with Betaseries api
func NewBetaClient(apiKey, login, password string) *BetaClient {

	// Force Go Routine to set token.
	finished := make(chan bool)

	var bs = &BetaClient{
		APIKey: apiKey,
		Token:  " ",
	}

	go bs.getAuthToken(login, password, finished)
	<-finished

	return bs
}

func (bs *BetaClient) getAuthToken(login, password string, finished chan bool) {
	hasher := md5.New()
	hasher.Write([]byte(password))

	var url = baseURL + "members/auth"
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
	data.Set("client_id", bs.APIKey)

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

// ***************************************************
//
//					EPISODE PART
//
// ***************************************************

// GetListEpisode return unWatched Episodes of all Show
// Important parameter, other parameter are on the [string]string other ( Check Documentation to learn more about it )
func (bs *BetaClient) GetListEpisode(listOfShowsID []string, typeOfShowID string, other map[string]string) string {

	// Definition type of show id, to rename request parameter
	var showID string
	if typeOfShowID == "TheTVDB" {
		showID = "showTheTVDBId"
	} else if typeOfShowID == "IMDB" {
		showID = "showIMDBId"
	} else {
		showID = "showId"
	}

	var url = baseURL + "episodes/list"
	// Setting request parameter
	var params = map[string]string{
		"token": bs.Token,
		showID:  strings.Join(listOfShowsID, ","),
	}

	// If other parameter set by user, adding them to params value.
	if len(other) > 0 {
		for keyOther, valueOther := range other {
			params[keyOther] = valueOther
		}
	}

	result := bs.makeRequest(url, "GET", params)
	return result
}

// GetLastEpisodeForShow return last Episodes of show(s)
// Important parameter, other parameter are on the [string]string other ( Check Documentation to learn more about it )
func (bs *BetaClient) GetLastEpisodeForShow(listOfShowsID []string, typeOfShowID string, displaySpecial bool) string {

	// Definition type of show id, to rename request parameter
	var showID string
	if typeOfShowID == "TheTVDB" {
		showID = "showTheTVDBId"
	} else {
		showID = "id"
	}

	var url = baseURL + "episodes/latest"
	// Definition request parameter
	var params = map[string]string{
		"token":    bs.Token,
		showID:     strings.Join(listOfShowsID, ","),
		"specials": strconv.FormatBool(displaySpecial),
	}
	result := bs.makeRequest(url, "GET", params)
	return result
}

// WatchedEpisode update betaserie's episode to passe hime to watched episode.
func (bs *BetaClient) WatchedEpisode(listOfShowsID []string, typeOfEpisodeID string, bulk bool, delete bool, note int) string {
	// Definition type of show id, to rename request parameter
	var showID string
	if typeOfEpisodeID == "TheTVDB" {
		showID = "showTheTVDBId"
	} else {
		showID = "id"
	}

	// note need to be between 1 and 5
	if note < 0 {
		note = 1
	} else if note > 5 {
		note = 5
	}

	var url = baseURL + "episodes/watched"
	// Set Request parameter like Api want...
	// Need to converte all parameter to String
	var params = map[string]string{
		"token":  bs.Token,
		showID:   strings.Join(listOfShowsID, ","),
		"bulk":   strconv.FormatBool(bulk),
		"delete": strconv.FormatBool(delete),
		"note":   strconv.FormatInt(int64(note), 10),
	}
	result := bs.makeRequest(url, "POST", params)
	return result

}

// UnWatched Episode update betaserie's episode to passe hime to unwatched episode.
func (bs *BetaClient) UnWatched(listOfShowsID []string, typeOfEpisodeID string) string {
	// Definition type of show id, to rename request parameter
	var showID string
	if typeOfEpisodeID == "TheTVDB" {
		showID = "showTheTVDBId"
	} else {
		showID = "id"
	}

	var url = baseURL + "episodes/watched"
	var params = map[string]string{
		"token": bs.Token,
		showID:  strings.Join(listOfShowsID, ","),
	}
	result := bs.makeRequest(url, "DELETE", params)
	return result
}

// ***************************************************
//
//					COMMENT PART
//
// ***************************************************
