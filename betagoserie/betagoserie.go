package betagoserie

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
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
	apiKey string
	token  string
}

func NewBetaClient(apiKey, login, password string) *BetaClient {

	var bs = &BetaClient{
		apiKey: apiKey,
		token:  " ",
	}

	bs.getAuthToken(login, password)

	return bs
}

func (*BetaClient) getListEpisode() {
	// var url = baseUrl + "/episodes/list"
}

func (bs *BetaClient) getAuthToken(login, password string) {
	println("AuthToken")
	// if bs.token != "" {
	// 	return bs.token
	// }

	hasher := md5.New()
	hasher.Write([]byte(password))

	var url = baseUrl + "members/auth"
	var params = map[string]string{
		"login":    login,
		"password": hex.EncodeToString(hasher.Sum(nil)),
	}

	bs.makeRequest(url, "post", params)

}

func (bs *BetaClient) makeRequest(url, urlType string, params map[string]string) {

	println("makeRequest")

	// parameters := url_net.Values{}
	// if len(params) > 0 {
	// 	println(params)
	// 	for paramKey, paramValue := range params {
	// 		parameters.Add(paramKey, paramValue)
	// 	}
	// }
	//
	// parameters.Add("client_id", bs.apiKey)

	// Build the request
	req, err := http.NewRequest(urlType, url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

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

	println(string(body))
}
