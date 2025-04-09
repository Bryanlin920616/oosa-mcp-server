package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// 景點結構
type Attraction struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Rating      float64  `json:"rating"`
	VisitHours  string   `json:"visitHours"`
	Tickets     Tickets  `json:"tickets"`
	Images      []string `json:"images"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
}

type Tickets struct {
	Adult  float64 `json:"adult"`
	Child  float64 `json:"child"`
	Senior float64 `json:"senior"`
}

// 全局變量
var attractions []Attraction

// 加載景點數據
func loadAttractions() error {
	file, err := os.Open("/Users/linjunli/Desktop/code/oosa-mcp-server/data/attractions.json")
	if err != nil {
		return fmt.Errorf("無法打開景點數據文件: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("無法讀取景點數據: %w", err)
	}

	if err := json.Unmarshal(data, &attractions); err != nil {
		return fmt.Errorf("無法解析景點數據: %w", err)
	}

	return nil
}

// 查找景點
func findAttractionByID(id string) *Attraction {
	for _, a := range attractions {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

// 搜索景點
func searchAttractions(query string) []Attraction {
	query = strings.ToLower(query)
	var results []Attraction

	for _, a := range attractions {
		if strings.Contains(strings.ToLower(a.Name), query) ||
			strings.Contains(strings.ToLower(a.Description), query) ||
			strings.Contains(strings.ToLower(a.Category), query) ||
			strings.Contains(strings.ToLower(a.Location), query) {
			results = append(results, a)
		}
	}

	return results
}

// 按類別過濾景點
func filterByCategory(category string) []Attraction {
	category = strings.ToLower(category)
	var results []Attraction

	for _, a := range attractions {
		if strings.ToLower(a.Category) == category {
			results = append(results, a)
		}
	}

	return results
}

// 資源處理函數 - 所有景點
func getAllAttractionsHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	data, err := json.Marshal(attractions)
	if err != nil {
		return nil, fmt.Errorf("無法序列化景點數據: %w", err)
	}

	// 使用NewReadResourceResult來建立回應
	result := mcp.NewReadResourceResult(string(data))
	return result.Contents, nil
}

// 資源處理函數 - 根據ID獲取景點
func getAttractionByIDHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// 從URI中解析ID
	parts := strings.Split(request.Params.URI, "/")
	id := parts[len(parts)-1]

	attraction := findAttractionByID(id)
	if attraction == nil {
		return nil, fmt.Errorf("找不到ID為 %s 的景點", id)
	}

	data, err := json.Marshal(attraction)
	if err != nil {
		return nil, fmt.Errorf("無法序列化景點數據: %w", err)
	}

	// 使用NewReadResourceResult來建立回應
	result := mcp.NewReadResourceResult(string(data))
	return result.Contents, nil
}

// 工具處理函數 - 搜索景點
func searchAttractionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var query string

	if q, ok := request.Params.Arguments["query"]; ok {
		query, _ = q.(string)
	}

	if query == "" {
		return mcp.NewToolResultError("查詢不能為空"), nil
	}

	results := searchAttractions(query)
	data, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("無法序列化搜索結果: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// 工具處理函數 - 按類別過濾景點
func filterByCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var category string

	if c, ok := request.Params.Arguments["category"]; ok {
		category, _ = c.(string)
	}

	if category == "" {
		return mcp.NewToolResultError("類別不能為空"), nil
	}

	results := filterByCategory(category)
	data, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("無法序列化過濾結果: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

func main() {
	// 加載景點數據
	if err := loadAttractions(); err != nil {
		log.Fatalf("無法加載景點數據: %v", err)
	}

	// 創建MCP服務器
	s := server.NewMCPServer(
		"OOSA 景點服務",
		"1.0.0",
		server.WithLogging(),
	)

	// 註冊資源
	s.AddResource(
		mcp.NewResource("attractions", "OOSA景點列表",
			mcp.WithResourceDescription("獲取所有OOSA景點數據")),
		getAllAttractionsHandler,
	)

	s.AddResource(
		mcp.NewResource("attraction/{id}", "根據ID獲取景點",
			mcp.WithResourceDescription("根據ID獲取特定景點的詳細信息")),
		getAttractionByIDHandler,
	)

	// 註冊工具 - 簡化版本，減少自定義參數
	s.AddTool(
		mcp.NewTool("search_attractions"),
		searchAttractionsHandler,
	)

	s.AddTool(
		mcp.NewTool("filter_by_category"),
		filterByCategoryHandler,
	)

	// 啟動服務器
	log.Println("OOSA 景點服務器已啟動...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("服務器錯誤: %v", err)
	}
}
