package main

import (
	"context"
)

func main() {
	parseFlags()
	err := run(context.Background())
	if err != nil {
		panic(err)
	}
}
