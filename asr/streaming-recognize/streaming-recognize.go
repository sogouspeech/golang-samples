// Copyright 2018 Sogou Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command streaming-recognize pipes the captured audio data to
// Sogou Speech API and outputs the transcript.
package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gordonklaus/portaudio"
	"github.com/modern-go/concurrent"
	asrv1 "golang.speech.sogou.com/apis/asr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const usage = `Usage: streaming-recognize
env SOGOU_SPEECH_ENDPOINT, SOGOU_SPEECH_APPID, SOGOU_SPEECH_TOKEN
must be set.
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

	if SogouSpeechEndpoint == "" || SogouSpeechAppID == "" || SogouSpeechToken == "" {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	fmt.Println("Start streaming-recognize.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	portaudio.Initialize()
	defer portaudio.Terminate()

	r, w := io.Pipe()

	executor := concurrent.NewUnboundedExecutor()

	executor.Go(func(ctx context.Context) { captureAudio(ctx, w) })
	executor.Go(func(ctx context.Context) { streamingRecognize(ctx, r) })

	<-sig
	fmt.Println("Exiting.")
	executor.StopAndWaitForever()
}

func captureAudio(ctx context.Context, w io.WriteCloser) {

	defer w.Close()

	in := make([]int16, 1600)
	audioStream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		panic(err)
	}

	defer audioStream.Close()

	if err = audioStream.Start(); err != nil {
		log.Fatal(err)
	}

	for {

		if err := audioStream.Read(); err != nil {
			log.Fatal(err)
		}

		if err := binary.Write(w, binary.LittleEndian, in); err != nil {
			log.Fatal(err)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func streamingRecognize(_ context.Context, r io.Reader) {

	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"appid", SogouSpeechAppID,
		"authorization", "Bearer "+SogouSpeechToken)

	dialOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))}
	conn, err := grpc.Dial(SogouSpeechEndpoint, dialOpts...)
	if err != nil {
		log.Fatal(err)
	}

	client := asrv1.NewAsrClient(conn)
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&asrv1.StreamingRecognizeRequest{
		StreamingRequest: &asrv1.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &asrv1.StreamingRecognitionConfig{
				Config: &asrv1.RecognitionConfig{
					Encoding:        asrv1.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "zh-cmn-Hans-CN",
				},
				InterimResults: true,
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		buf := make([]byte, 3200)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				if err := stream.Send(&asrv1.StreamingRecognizeRequest{
					StreamingRequest: &asrv1.StreamingRecognizeRequest_AudioContent{
						AudioContent: buf[:n],
					},
				}); err != nil {
					log.Printf("Could not send audio: %v", err)
				}
			}
			if err == io.EOF {
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Could not close stream: %v", err)
				}
				return
			}
			if err != nil {
				log.Printf("Could not read audio: %v", err)
				continue
			}
		}
	}()

	m := jsonpb.Marshaler{OrigName: true}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot stream results: %v", err)
		}
		if err := resp.GetError(); err != nil {
			log.Fatalf("Could not recognize: %v", err)
		}

		res, _ := m.MarshalToString(resp)
		fmt.Println(res)
	}

}
