package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	pluginVersion = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-cobertura-s3"
	app.Usage = "Drone plugin to upload cobertura code coverage reports to AWS S3 bucket and publish the s3 bucket static site url under 'Executions > Artifacts' tab"
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

	if roleArn == "" {
		fmt.Printf("\nAWS configuration....")
		exec.Command("aws", "configure", "set", "aws_access_key_id", awsAccessKey).Run()
		exec.Command("aws", "configure", "set", "aws_secret_access_key", awsSecretKey).Run()
	} else {
		if roleSessionName == "" {
			roleSessionName = "drone"
		}
		cmd := exec.Command("/bin/sh", "-c",
			`export $(printf "AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s AWS_SESSION_TOKEN=%s" $(aws sts assume-role-with-web-identity --role-arn `+roleArn+` --role-session-name `+roleSessionName+`--web-identity-token file://$AWS_WEB_IDENTITY_TOKEN_FILE --query "Credentials.[AccessKeyId,SecretAccessKey,SessionToken]" --output text))`)

		if err := cmd.Run(); err != nil {
			log.Fatal("Error:", err)
		}
	}

	reportUploadcmd := exec.Command("aws", "s3", "cp", reportSource, "s3://"+awsBucket+"/"+newFolder, "--region", awsDefaultRegion, "--recursive", "--acl=public-read")
	fmt.Printf("\nreportUploadcmd: %v", reportUploadcmd)

	out, err := reportUploadcmd.Output()
	if err != nil {
		// Error here
		fmt.Printf("\nAWS CLI Error: %s\n", reportUploadcmd.Stderr)
		return err
	}
	fmt.Printf("\nOutput: %s\n", out)

	urls := "http://" + awsBucket + ".s3-website." + awsDefaultRegion + ".amazonaws.com/" + newFolder + "/index.html"
	artifactFilePath := c.String("artifact-file")

	files := make([]File, 0)
	files = append(files, File{Name: artifactFilePath, URL: urls})

	return writeArtifactFile(files, artifactFilePath)
}
