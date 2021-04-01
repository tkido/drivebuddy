package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	transcoder "github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/polly"

	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	doTranscode()
}

const (
	region          = `ap-northeast-1`
	presetId        = `1351620000001-200060`
	pipelineId      = `1614661200392-wu2bfw`
	inputKey        = `glass.mp3`
	outputKey       = `glass`
	outputKeyPrefix = `awssdk/`
	segmentDuration = `15`
)

func doTranscode() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := transcoder.New(sess)

	resp, err := svc.CreateJob(
		&transcoder.CreateJobInput{
			Input: &transcoder.JobInput{
				Key: aws.String(inputKey),
			},
			OutputKeyPrefix: aws.String(outputKeyPrefix),
			Outputs: []*transcoder.CreateJobOutput{
				&transcoder.CreateJobOutput{
					PresetId:        aws.String(presetId),
					Key:             aws.String(outputKey),
					SegmentDuration: aws.String(segmentDuration),
				},
			},
			PipelineId: aws.String(pipelineId),
		},
	)
	if err != nil {
		log.Printf("Failed: Create Job, %v\n", err)
		return
	}

	log.Printf("Job Response: %v\n", resp.Job)
}

func doPolly() {
	// The name of the text file to convert to MP3
	fileName := `_testdata/glass.txt`

	// Open text file and get it's contents as a string
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Got error opening file " + fileName)
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Convert bytes to string
	s := string(contents[:])

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create Polly client
	svc := polly.New(sess)

	// Output to MP3 using voice Joanna
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Joanna")}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("Got error calling SynthesizeSpeech:")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Save as MP3
	names := strings.Split(fileName, ".")
	name := names[0]
	mp3File := name + ".mp3"

	outFile, err := os.Create(mp3File)
	if err != nil {
		fmt.Println("Got error creating " + mp3File + ":")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		fmt.Println("Got error saving MP3:")
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
