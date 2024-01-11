# drone-cobertura-s3

Drone plugin to upload cobertura code coverage reports to AWS S3 bucket and publish the bucket static site URL to `Artifacts` tab under `Pipieline Execution`.

## Build

Build the binary for different OS/Arch with the following commands:

```bash
env GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/drone-cobertura-s3 .
env GOOS=linux GOARCH=arm64 go build -o release/linux/arm64/drone-cobertura-s3 .
env GOOS=windows go build -o release/windows/amd64/drone-cobertura-s3.exe   
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
