package main

import (
	"context"
	"log"
	"strconv"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/common/helpers"
	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/presentation"
)

func main() {
	ctx := context.Background()

	port, err := strconv.Atoi(helpers.MustGetEnvVar("PORT"))
	if err != nil {
		log.Panicf("Could not get environment variable %s", helpers.MustGetEnvVar("PORT"))
	}

	presentation.PrepareServer(ctx, port)
}
