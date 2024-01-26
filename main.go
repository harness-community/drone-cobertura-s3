package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	pluginVersion = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-cobertura-s3"
	app.Usage = "Drone plugin to upload cobertura code coverage reports to AWS S3 bucket and publish the s3 bucket static site URL under 'Executions > Artifacts' tab"
	app.Action = run
	app.Version = pluginVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key ID",
			EnvVar: "PLUGIN_AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Access Key",
			EnvVar: "PLUGIN_AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-default-region",
			Usage:  "AWS Default Region",
			EnvVar: "PLUGIN_AWS_DEFAULT_REGION",
		},
		cli.StringFlag{
			Name:   "aws-bucket",
			Usage:  "The name of the S3 bucket where you want to upload the Cobertura report files",
			EnvVar: "PLUGIN_AWS_BUCKET",
		},
		cli.StringFlag{
			Name:   "role-arn",
			Usage:  "AWS Role Arn",
			EnvVar: "PLUGIN_ROLE_ARN",
		},
		cli.StringFlag{
			Name:   "role-session-name",
			Usage:  "AWS Role Session Name",
			EnvVar: "PLUGIN_ROLE_SESSION_NAME",
		},
		cli.StringFlag{
			Name:   "report-source",
			Usage:  "The source directory containing all the Cobertura report files which are to be uploaded to S3 bucket",
			EnvVar: "PLUGIN_REPORT_SOURCE",
		},
		cli.StringFlag{
			Name:   "report-target",
			Usage:  "The target directory in the S3 bucket where the user wants to upload the Cobertura report files",
			EnvVar: "PLUGIN_REPORT_TARGET",
		},
		cli.StringFlag{
			Name:   "artifact-file",
			Usage:  "Artifact file",
			EnvVar: "PLUGIN_ARTIFACT_FILE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	pipelineSeqID := os.Getenv("DRONE_BUILD_NUMBER")
	awsAccessKey := c.String("aws-access-key")
	awsSecretKey := c.String("aws-secret-key")
	awsDefaultRegion := c.String("aws-default-region")
	awsBucket := c.String("aws-bucket")
	roleArn := c.String("role-arn")
	roleSessionName := c.String("role-session-name")
	reportSource := c.String("report-source")
	reportTarget := c.String("report-target")

	var newFolder string
	if reportTarget == "" {
		newFolder = "build-" + pipelineSeqID
	} else {
		newFolder = reportTarget + "/build-" + pipelineSeqID
	}

	fmt.Printf("\nUploading Cobertura reports to " + awsBucket + "/" + newFolder)

	// AWS session and S3 client setup
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsDefaultRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}))

	if roleArn != "" {
		// Assume Role if roleArn is provided
		if roleSessionName == "" {
			roleSessionName = "drone"
		}
		creds, err := stsAssumeRole(sess, roleArn, roleSessionName)
		if err != nil {
			fmt.Printf("\nError: %s\n", err)
			return err
		}
		sess.Config.Credentials = creds
	}

	// Check if reportSource is a directory
	fileInfo, err := os.Stat(reportSource)
	if err != nil {
		fmt.Printf("\nError checking file or directory: %s\n", err)
		return err
	}

	if !fileInfo.IsDir() {
		fmt.Printf("\nError: The specified 'report-source' is not a directory. Please provide a valid directory path.\n")
		return fmt.Errorf("report-source is not a directory")
	}

	// Walk through the directory and upload files
	err = filepath.Walk(reportSource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("\nError traversing directory: %s\n", err)
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("\nError opening file: %s\n", err)
			return err
		}
		defer file.Close()

		// Construct the S3 key relative to the original reportSource directory
		key := filepath.ToSlash(filepath.Join(newFolder, strings.TrimPrefix(path, reportSource)))

		// Upload report to S3
		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(awsBucket),
			Key:    aws.String(key),
			Body:   file,
			ACL:    aws.String("public-read"),
		})
		if err != nil {
			fmt.Printf("\nAWS SDK Error: %s\n", err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("\nError walking through directory: %s\n", err)
		return err
	}

	urls := fmt.Sprintf("http://%s.s3-website.%s.amazonaws.com/%s/index.html", awsBucket, awsDefaultRegion, newFolder)
	artifactFilePath := c.String("artifact-file")

	files := make([]File, 0)
	files = append(files, File{Name: artifactFilePath, URL: urls})

	return writeArtifactFile(files, artifactFilePath)
}

// stsAssumeRole assumes a role and returns the temporary credentials
func stsAssumeRole(sess *session.Session, roleArn, roleSessionName string) (*credentials.Credentials, error) {
	svc := sts.New(sess)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(roleSessionName),
	}

	result, err := svc.AssumeRole(input)
	if err != nil {
		return nil, err
	}

	return credentials.NewStaticCredentials(
		*result.Credentials.AccessKeyId,
		*result.Credentials.SecretAccessKey,
		*result.Credentials.SessionToken,
	), nil
}
