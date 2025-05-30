package tools

import (
	"context"
	"encoding/json"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func retryCloudBuildJob() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("retry_cloud_build_job",
			mcp.WithDescription("Retry a failed or cancelled Google Cloud Build job for a project"),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("GCP project ID"),
			),
			mcp.WithString("build_id",
				mcp.Required(),
				mcp.Description("Cloud Build job ID to retry"),
			),
		),
		Handler: retryCloudBuildJobHandler,
	}
}

func retryCloudBuildJobHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError("Missing project_id: " + err.Error()), nil
	}
	buildID, err := request.RequireString("build_id")
	if err != nil {
		return mcp.NewToolResultError("Missing build_id: " + err.Error()), nil
	}

	client, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return mcp.NewToolResultError("Failed to create Cloud Build client: " + err.Error()), nil
	}
	defer client.Close()

	operation, err := client.RetryBuild(ctx, &cloudbuildpb.RetryBuildRequest{
		ProjectId: projectID,
		Id:        buildID,
	})
	if err != nil {
		return mcp.NewToolResultError("Error retrying build: " + err.Error()), nil
	}

	operationJSON, _ := json.Marshal(operation)
	return mcp.NewToolResultText(string(operationJSON)), nil
}
