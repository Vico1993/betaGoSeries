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

// TODO : handle
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

type errorStruct struct {
}

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

	var url = baseURL + "episodes/list"
	var showID = getParameterOfShowType(typeOfShowID)

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

	var url = baseURL + "episodes/latest"
	var showID = getParameterOfShowType(typeOfShowID)

	// Definition request parameter
	var params = map[string]string{
		"token":    bs.Token,
		showID:     strings.Join(listOfShowsID, ","),
		"specials": strconv.FormatBool(displaySpecial),
	}
	result := bs.makeRequest(url, "GET", params)
	return result
}

// GetEpisodeByFileName return betaseries episode by is file name
func (bs *BetaClient) GetEpisodeByFileName(filename string) string {
	var url = baseURL + "episodes/scrapper"

	var params = map[string]string{
		"token": bs.Token,
		"file":  filename,
	}
	result := bs.makeRequest(url, "GET", params)
	return result
}

// WatchedEpisode update betaserie's episode to passe hime to watched episode.
func (bs *BetaClient) WatchedEpisode(listOfShowsID []string, typeOfEpisodeID string, bulk bool, delete bool, note int) string {
	// note need to be between 1 and 5
	if note < 0 {
		note = 1
	} else if note > 5 {
		note = 5
	}

	var url = baseURL + "episodes/watched"
	var showID = getParameterOfShowType(typeOfEpisodeID)

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

	var url = baseURL + "episodes/watched"
	var showID = getParameterOfShowType(typeOfEpisodeID)

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

// ***************************************************
//
//					Tool Client's Needed
//
// ***************************************************

func getParameterOfShowType(typeOfShowID string) string {
	// Definition type of show id, to rename request parameter
	var showID string
	if typeOfShowID == "TheTVDB" {
		showID = "showTheTVDBId"
	} else if typeOfShowID == "IMDB" {
		showID = "showIMDBId"
	} else {
		showID = "showId"
	}

	return showID
}
