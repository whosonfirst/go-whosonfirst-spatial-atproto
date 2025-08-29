package main

import (
	"context"
	"log"

	_ "github.com/whosonfirst/go-reader-cachereader/v2"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-duckdb"

	"github.com/whosonfirst/go-whosonfirst-spatial-atproto/app/server"
)

func main() {

	ctx := context.Background()
	err := server.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
