package tools

import (
	"context"
	"encoding/json"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getCloudBuildJob() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_cloud_build_job",
			mcp.WithDescription("Get a specific Google Cloud Build job by ID for a project"),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("GCP project ID"),
			),
			mcp.WithString("build_id",
				mcp.Required(),
				mcp.Description("Cloud Build job ID"),
			),
		),
		Handler: getCloudBuildJobHandler,
	}
}

func getCloudBuildJobHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	build, err := client.GetBuild(ctx, &cloudbuildpb.GetBuildRequest{
		ProjectId: projectID,
		Id:        buildID,
	})
	if err != nil {
		return mcp.NewToolResultError("Error getting build: " + err.Error()), nil
	}

	buildInfo := map[string]interface{}{
		"id":          build.Id,
		"status":      build.Status.String(),
		"create_time": build.CreateTime,
		"start_time":  build.StartTime,
		"finish_time": build.FinishTime,
		"steps":       build.Steps,
		"images":      build.Images,
	}
	resultBytes, _ := json.Marshal(buildInfo)
	return mcp.NewToolResultText(string(resultBytes)), nil
}
