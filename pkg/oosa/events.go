package oosa

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetEvents(client *Client) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_events",
			mcp.WithString("event_past",
				mcp.Description("Filter events that occurred in the past"),
			),
			mcp.WithString("event_period_begin",
				mcp.Description("The beginning of the event period"),
			),
			mcp.WithString("event_period_end",
				mcp.Description("The end of the event period"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			eventPast, err := OptionalParam[string](request, "event_past")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			eventPeriodBegin, err := OptionalParam[string](request, "event_period_begin")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			eventPeriodEnd, err := OptionalParam[string](request, "event_period_end")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// TODO: use real mongo client to get events
			events, err := client.GetEvents(ctx, eventPast, eventPeriodBegin, eventPeriodEnd)
			if err != nil {
				return nil, fmt.Errorf("failed to get events: %w", err)
			}

			// TODO: StatusCode check

			r, err := json.Marshal(events)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal events: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
