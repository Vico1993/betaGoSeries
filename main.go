package main

import (
	"github.com/vico1993/betaGoSeries/betagoserie"
)

func main() {
	// Connection
	var client = betagoserie.NewBetaClient( #apiKey#, #login#, #password# )

	// test := client.GetListEpisode()

	test := client.GetLastEpisodeForShow([]string{"11001"}, "betaseries_id", false)
	println(test)
}
