// Copyright 2018 Sogou Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command recognize sends audio data to the Sogou Speech API
// and prints its transcript.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	asrv1 "golang.speech.sogou.com/apis/asr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const usage = `Usage: recognize <audiofile>
Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.
`

var (
	SogouSpeechEndpoint = ""
	SogouSpeechAppID    = ""
	SogouSpeechToken    = ""
)

func init() {

	SogouSpeechEndpoint = os.Getenv("SOGOU_SPEECH_ENDPOINT")
	SogouSpeechAppID = os.Getenv("SOGOU_SPEECH_APPID")
	SogouSpeechToken = os.Getenv("SOGOU_SPEECH_TOKEN")

}

func main() {

	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	// Perform the request.
	if err := recognize(os.Stdout, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func recognize(w io.Writer, file string) error {

	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"appid", SogouSpeechAppID,
		"authorization", "Bearer " + SogouSpeechToken)

	dialOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))}
	conn, err := grpc.Dial(SogouSpeechEndpoint, dialOpts...)
	if err != nil {
		return err
	}

	client := asrv1.NewAsrClient(conn)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &asrv1.RecognizeRequest{
		Config: &asrv1.RecognitionConfig{
			Encoding:        asrv1.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "zh-cmn-Hans-CN",
		},
		Audio: &asrv1.RecognitionAudio{
			AudioSource: &asrv1.RecognitionAudio_Content{Content: data},
		},
	})

	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.GetResults() {
		for _, alt := range result.GetAlternatives() {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.GetTranscript(), alt.GetConfidence())
		}
	}
	return nil
}