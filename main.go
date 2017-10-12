package main

import (
	"github.com/vico1993/betaGoSeries/betagoserie"
)

func main() {
	// Connection
	var client = betagoserie.NewBetaClient("ee7422ce11a2", "Vico1993", "victor1993")

	// test := client.GetListEpisode()

	test := client.GetLastEpisodeForShow([]string{"11001"}, "betaseries_id", false)
	println(test)
}
