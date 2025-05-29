package tools

import (
	"context"
	"encoding/json"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/api/iterator"
)

func listCloudBuildJobs() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_cloud_build_jobs",
			mcp.WithDescription("List Google Cloud Build jobs for a project"),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("GCP project ID to list builds for"),
			),
		),
		Handler: listCloudBuildJobsHandler,
	}
}

func listCloudBuildJobsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError("Missing project_id: " + err.Error()), nil
	}

	client, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return mcp.NewToolResultError("Failed to create Cloud Build client: " + err.Error()), nil
	}
	defer client.Close()

	it := client.ListBuilds(ctx, &cloudbuildpb.ListBuildsRequest{ProjectId: projectID})
	var jobs []map[string]interface{}
	for {
		build, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return mcp.NewToolResultError("Error listing builds: " + err.Error()), nil
		}
		jobs = append(jobs, map[string]interface{}{
			"id":          build.Id,
			"status":      build.Status.String(),
			"create_time": build.CreateTime,
			"start_time":  build.StartTime,
			"finish_time": build.FinishTime,
		})
	}
	result := map[string]interface{}{
		"project_id": projectID,
		"jobs":       jobs,
	}
	resultBytes, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultBytes)), nil
}
