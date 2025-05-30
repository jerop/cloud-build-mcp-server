# Cloud Build MCP Server

A Model Context Protocol (MCP) server for interacting with Google Cloud Build.

## Tools

| Feature                                   | Tool                | Parameters                                 |
|--------------------------------------------|--------------------------|---------------------------------------------|
| List Cloud Build jobs for a project        | `list_cloud_build_jobs`  | `project_id`                                |
| View details of a specific Cloud Build job | `get_cloud_build_job`    | `project_id`, `build_id`                    |
| Create new Cloud Build jobs                | `create_cloud_build_job` | `project_id`, `build_config_json`           |
| Retry failed or cancelled Cloud Build jobs | `retry_cloud_build_job`  | `project_id`, `build_id`                    |

## Prerequisites

- Go 1.21+
- Google Cloud project with Cloud Build API enabled
- [Google Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/provide-credentials-adc):
  ```sh
  gcloud auth application-default login
  ```

## Build and Run

### 1. Build the Server

```sh
cd cmd
go build -o cloud-build-mcp-server
```

### 2. Run the Server

```sh
./cloud-build-mcp-server
```

## Use in Agent

TODO(jerop): add instructions

