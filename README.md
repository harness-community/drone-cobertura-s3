# drone-cobertura-s3

- [Synopsis](#Synopsis)
- [Parameters](#Parameters)
- [Notes](#Notes)
- [Plugin Image](#Plugin-Image)
- [Examples](#Examples)

## Synopsis
This plugin is designed to streamline the process of uploading Cobertura reports to an Amazon S3 bucket and retrieving and publishing a URL to `Artifacts` tab under `Pipieline Execution` for users to view their reports. With this plugin, users can easily upload their Cobertura reports along with a custom HTML template to S3, and obtain a URL pointing to their customized report visualization.

## Parameters
| Parameter                                                                                                                              | Comments                                                        |
| :------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| aws_access_key_id <span style="font-size: 10px"><br/>`string`</span> <span style="color:red; font-size: 10px">`required`</span>               |  The unique ID to access your AWS account and AWS CLI.                |
| aws_secret_access_key <span style="font-size: 10px"><br/>`string`</span> <span style="color:red; font-size: 10px">`required`</span>                  |  The security token to connect to AWS.                            |
| aws_default_region <span style="font-size: 10px"><br/>`string`</span> <span style="color:red; font-size: 10px">`required`</span>              |      The default region where the bucket is deployed.                        |
| aws_bucket <span style="font-size: 10px"><br/>`string`</span> <span style="color:red; font-size: 10px">`required`</span> |                   The name of the S3 bucket where the files are to be uploaded.                       |
| report_source <span style="font-size: 10px"><br/>`string`</span><span style="color:red; font-size: 10px">`required`</span>                     | The directory containing the Cobertura report files and custom HTML template. |
| role_arn <span style="font-size: 10px"><br/>`string`</span>                   | The directory containing the Cobertura report files and custom HTML template. |
| role_session_name <span style="font-size: 10px"><br/>`string`</span>                   | The directory containing the Cobertura report files and custom HTML template. |
| report_target <span style="font-size: 10px"><br/>`string`</span>                   | The target directory in the S3 bucket where the Cobertura report files is to be uploaded. |
| artifact_file <span style="font-size: 10px"><br/>`string`</span>                   | The directory containing the Cobertura report files and custom HTML template. |


## Notes
Users are required to create a custom HTML template named `index.html` to define the design and layout of their Cobertura report visualization. The `index.html` must be stored in `report_source`. This HTML file will serve as the frontend interface for viewing the uploaded reports. It should include placeholders or dynamic elements where relevant information from the Cobertura reports can be displayed.

## Plugin Image

The plugin `harnesscommunitytest/drone-cobertura-s3` is available for the following architectures:

| OS            | Tag                                |
| ------------- | ---------------------------------- |
| latest        | `linux-amd64/arm64, windows-amd64` |
| linux/amd64   | `linux-amd64`                      |
| linux/arm64   | `linux-arm64`                      |
| windows/amd64 | `windows-amd64`                    |


## Examples

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
                      aws_bucket: my-bucket
                      report_source: /harness
```
