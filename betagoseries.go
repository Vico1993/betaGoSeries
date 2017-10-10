package betagoserie

// TODO : put const into config file ( ENV VARIABLE ? )
// const apiKey = "ee7422ce11a2"
// const apiSecret = "d2e555a996fb64e49febb5adb9d1c818"

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

type betaClient struct {
	apiKey    string
	apiSecret string
}

func (*betaClient) getListEpisode() {
	// var url = baseUrl + "/episodes/list"
}

func (*betaClient) makeRequest(url, urlType string, params []string) {

}
