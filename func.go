package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	config, err := readConfig(in)
	if err != nil {
		outputError(out, err)
		return
	}
	output, err := bootstrap(ctx, config)
	if err != nil {
		outputError(out, err)
		return
	}
	json.NewEncoder(out).Encode(output)
}

func outputError(out io.Writer, err error) {
	msg := struct {
		Msg string `json:"error"`
	}{
		Msg: fmt.Sprintf("failed to bootstrap due to: %s", err),
	}
	json.NewEncoder(out).Encode(&msg)
}
