package tools

import (
	"context"
	"encoding/json"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func createCloudBuildJob() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("create_cloud_build_job",
			mcp.WithDescription("Create a new Google Cloud Build job for a project"),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("GCP project ID to create the build in"),
			),
			mcp.WithString("build_config_json",
				mcp.Required(),
				mcp.Description("Build configuration as a JSON string (matches google.cloudbuild.v1.Build)"),
			),
		),
		Handler: createCloudBuildJobHandler,
	}
}

func createCloudBuildJobHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError("Missing project_id: " + err.Error()), nil
	}
	buildConfigJSON, err := request.RequireString("build_config_json")
	if err != nil {
		return mcp.NewToolResultError("Missing build_config_json: " + err.Error()), nil
	}

	var build cloudbuildpb.Build
	if err := json.Unmarshal([]byte(buildConfigJSON), &build); err != nil {
		return mcp.NewToolResultError("Invalid build_config_json: " + err.Error()), nil
	}

	client, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return mcp.NewToolResultError("Failed to create Cloud Build client: " + err.Error()), nil
	}
	defer client.Close()

	operation, err := client.CreateBuild(ctx, &cloudbuildpb.CreateBuildRequest{
		ProjectId: projectID,
		Build:     &build,
	})
	if err != nil {
		return mcp.NewToolResultError("Error creating build: " + err.Error()), nil
	}

	operationJSON, _ := json.Marshal(operation)
	return mcp.NewToolResultText(string(operationJSON)), nil
}
