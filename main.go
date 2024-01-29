package main

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	artifactFilePath := c.String("artifact-file")

	newFolder := getNewFolder(pipelineSeqID, reportTarget)

	fmt.Printf("\nUploading Cobertura reports to %s/%s", awsBucket, newFolder)

	creds := getAWSCredentials(awsAccessKey, awsSecretKey, roleArn, roleSessionName, awsDefaultRegion)

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsDefaultRegion),
		Credentials: creds,
	}))

	// Upload each file in the directory to S3
	err := filepath.Walk(reportSource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal("Error opening file:", err)
			return err
		}
		defer file.Close()

		// Calculate the S3 key based on the original file path
		s3Key := strings.TrimPrefix(path, reportSource)
		s3Key = strings.TrimLeft(s3Key, string(os.PathSeparator))
		s3Key = filepath.Join(newFolder, s3Key)

		params := &s3.PutObjectInput{
			Bucket:      aws.String(awsBucket),
			Key:         aws.String(s3Key),
			Body:        file,
			ACL:         aws.String("public-read"),
			ContentType: aws.String("text/html"),
		}

		s3Svc := s3.New(sess)
		_, err = s3Svc.PutObject(params)
		if err != nil {
			fmt.Printf("\nAWS SDK Error: %v\n", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error walking directory:", err)
		return err
	}

	urls := fmt.Sprintf("http://%s.s3-website.%s.amazonaws.com/%s/index.html", awsBucket, awsDefaultRegion, newFolder)

	files := make([]File, 0)
	files = append(files, File{Name: artifactFilePath, URL: urls})

	return writeArtifactFile(files, artifactFilePath)
}

func getNewFolder(pipelineSeqID, reportTarget string) string {
	if reportTarget == "" {
		return "build-" + pipelineSeqID
	}
	return reportTarget + "/build-" + pipelineSeqID
}

func getAWSCredentials(awsAccessKey, awsSecretKey, roleArn, roleSessionName, awsDefaultRegion string) *credentials.Credentials {
	if roleArn == "" {
		return credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		Region:      aws.String(awsDefaultRegion),
	}))

	stsSvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(roleSessionName),
	}

	resp, err := stsSvc.AssumeRole(params)
	if err != nil {
		log.Fatal("Error assuming role:", err)
	}

	return credentials.NewStaticCredentials(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken)
}
