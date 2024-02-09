# drone-cobertura-s3

- [Synopsis](#Synopsis)
- [Notes](#Notes)
- [Build](#Build)
- [Docker](#Docker)
- [Usage](#Usage)

## Synopsis
This plugin is designed to streamline the process of uploading Cobertura reports to an Amazon S3 bucket and retrieving and publishing a URL to `Artifacts` tab under `Pipieline Execution` for users to view their reports. With this plugin, users can easily upload their Cobertura reports along with a custom HTML template to S3, and obtain a URL pointing to their customized report visualization.

## Notes
Users are required to create a custom HTML template named index.html to define the design and layout of their Cobertura report visualization. This HTML file will serve as the frontend interface for viewing the uploaded reports. It should include placeholders or dynamic elements where relevant information from the Cobertura reports can be displayed.

## Build

Build the binary for different OS/Arch with the following command:

```bash
./scripts/build.sh
```

## Docker

Build the Docker image with the following commands:

```
docker build -t DOCKER_ORG/drone-cobertura-s3 -f PATH_TO_DOCKERFILE
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-cobertura-s3' not found or does not exist..
```

## Usage

```bash
docker run --rm \
  -e PLUGIN_AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
  -e PLUGIN_AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
  -e PLUGIN_AWS_DEFAULT_REGION=ap-southeast-2 \
  -e PLUGIN_AWS_BUCKET=bucket-name \
  -e PLUGIN_REPORT_SOURCE=/harness \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  harnesscommunitytest/drone-cobertura-s3
```



In Harness CI,
Using AWS Access Key ID and Secret Access Key as authentication method:
```yaml
              - step:
                  type: Plugin
                  name: drone_cobertura_s3
                  identifier: drone_cobertura_s3
                  spec:
                    connectorRef: account.harnessImage
                    image: harnesscommunitytest/drone-cobertura-s3
                    settings:
                      aws_access_key_id: <+secrets.getValue("awsaccesskeyid")>
                      aws_secret_access_key: <+secrets.getValue("awssecretaccesskey")>
                      aws_default_region: ap-south-1
                      aws_bucket: cobertura-plugin
                      report_source: /harness
```
