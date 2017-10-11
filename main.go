package main

import (
	"github.com/vico1993/betaGoSeries/betagoserie"
)

func main() {
	var client = betagoserie.NewBetaClient("ee7422ce11a2", "Vico1993", "victor1993")

	client.GetListEpisode()
}
