package oosa

import (
	"context"
	"time"
)

type Client struct {
}

type Event struct {
	ID               string              `json:"events_id"`
	Name             string              `json:"events_name"`
	Date             string              `json:"events_date"`
	DateEnd          string              `json:"events_date_end"`
	Deadline         string              `json:"events_deadline"`
	Place            string              `json:"events_place"`
	Lat              float64             `json:"events_lat"`
	Lng              float64             `json:"events_lng"`
	MeetingPointName string              `json:"events_meeting_point_name"`
	MeetingPointLat  float64             `json:"events_meeting_point_lat"`
	MeetingPointLng  float64             `json:"events_meeting_point_lng"`
	ParticipantLimit float64             `json:"events_participant_limit"`
	PaymentRequired  int64               `json:"events_payment_required"`
	PaymentFee       float64             `json:"events_payment_fee"`
	Photo            string              `json:"events_photo"`
	Type             string              `json:"events_type"`
	CreatedByUser    *UserAgg            `json:"events_created_by_user,omitempty"`
	Participants     *EventsParticipants `json:"events_participants,omitempty"`
	CreatedAt        *string             `json:"events_created_at,omitempty"`
}

type UserAgg struct {
	ID     string `json:"user_id"`
	Name   string `json:"user_name"`
	Email  string `json:"user_email"`
	Avatar string `json:"user_avatar"`
}

type EventsParticipants struct {
	LatestThreeUser []UserAgg `json:"latest_tree_user"`
	RemainNumber    int64     `json:"remain_number"`
}

type Idea struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
}

func (c *Client) GetEvents(ctx context.Context, eventPast string, eventPeriodBegin string, eventPeriodEnd string) ([]Event, error) {
	// 建立一些測試用戶
	users := []UserAgg{
		{
			ID:     "user1",
			Name:   "王小明",
			Email:  "xiaoming@example.com",
			Avatar: "https://example.com/avatar1.jpg",
		},
		{
			ID:     "user2",
			Name:   "李小華",
			Email:  "xiaohua@example.com",
			Avatar: "https://example.com/avatar2.jpg",
		},
		{
			ID:     "user3",
			Name:   "張大偉",
			Email:  "dawei@example.com",
			Avatar: "https://example.com/avatar3.jpg",
		},
	}

	// 建立活動列表
	events := []Event{
		{
			ID:               "event1",
			Name:             "象山步道健行賞夜景",
			Date:             "2025-04-11T16:00:00Z",
			DateEnd:          "2025-04-11T20:00:00Z",
			Deadline:         "2025-04-10T23:59:59Z",
			Place:            "象山步道",
			Lat:              25.0330,
			Lng:              121.5700,
			MeetingPointName: "象山捷運站2號出口",
			MeetingPointLat:  25.0330,
			MeetingPointLng:  121.5700,
			ParticipantLimit: 15,
			PaymentRequired:  0,
			PaymentFee:       0,
			Photo:            "https://example.com/elephant_mountain.jpg",
			Type:             "健行",
			CreatedByUser:    &users[0],
			Participants: &EventsParticipants{
				LatestThreeUser: users,
				RemainNumber:    12,
			},
			CreatedAt: strPtr("2025-03-11T12:00:00Z"),
		},
		{
			ID:               "event2",
			Name:             "大稻埕老街文化導覽",
			Date:             "2025-04-12T14:00:00Z",
			DateEnd:          "2025-04-12T18:00:00Z",
			Deadline:         "2025-04-11T23:59:59Z",
			Place:            "大稻埕",
			Lat:              25.0550,
			Lng:              121.5100,
			MeetingPointName: "大稻埕碼頭",
			MeetingPointLat:  25.0550,
			MeetingPointLng:  121.5100,
			ParticipantLimit: 20,
			PaymentRequired:  1,
			PaymentFee:       200,
			Photo:            "https://example.com/dadaocheng.jpg",
			Type:             "文化導覽",
			CreatedByUser:    &users[1],
			Participants: &EventsParticipants{
				LatestThreeUser: users[:2],
				RemainNumber:    17,
			},
			CreatedAt: strPtr("2025-03-12T12:00:00Z"),
		},
		{
			ID:               "event3",
			Name:             "北投溫泉泡湯之旅",
			Date:             "2025-04-13T13:00:00Z",
			DateEnd:          "2025-04-13T17:00:00Z",
			Deadline:         "2025-04-12T23:59:59Z",
			Place:            "北投溫泉博物館",
			Lat:              25.1370,
			Lng:              121.5070,
			MeetingPointName: "北投捷運站出口",
			MeetingPointLat:  25.1319,
			MeetingPointLng:  121.4986,
			ParticipantLimit: 10,
			PaymentRequired:  1,
			PaymentFee:       500,
			Photo:            "https://example.com/beitou_hot_spring.jpg",
			Type:             "溫泉",
			CreatedByUser:    &users[2],
			Participants: &EventsParticipants{
				LatestThreeUser: users[1:],
				RemainNumber:    7,
			},
			CreatedAt: strPtr("2025-03-13T12:00:00Z"),
		},
		{
			ID:               "event4",
			Name:             "士林夜市美食探索",
			Date:             "2025-04-14T18:00:00Z",
			DateEnd:          "2025-04-14T22:00:00Z",
			Deadline:         "2025-04-13T23:59:59Z",
			Place:            "士林夜市",
			Lat:              25.0880,
			Lng:              121.5200,
			MeetingPointName: "士林捷運站1號出口",
			MeetingPointLat:  25.0880,
			MeetingPointLng:  121.5200,
			ParticipantLimit: 8,
			PaymentRequired:  0,
			PaymentFee:       0,
			Photo:            "https://example.com/shilin_night_market.jpg",
			Type:             "美食",
			CreatedByUser:    &users[0],
			Participants: &EventsParticipants{
				LatestThreeUser: users[:2],
				RemainNumber:    5,
			},
			CreatedAt: strPtr("2025-03-14T12:00:00Z"),
		},
		{
			ID:               "event5",
			Name:             "陽明山賞花健行",
			Date:             "2025-04-15T09:00:00Z",
			DateEnd:          "2025-04-15T15:00:00Z",
			Deadline:         "2025-04-14T23:59:59Z",
			Place:            "陽明山國家公園",
			Lat:              25.1700,
			Lng:              121.5400,
			MeetingPointName: "劍潭捷運站1號出口",
			MeetingPointLat:  25.0836,
			MeetingPointLng:  121.5256,
			ParticipantLimit: 12,
			PaymentRequired:  1,
			PaymentFee:       300,
			Photo:            "https://example.com/yangmingshan.jpg",
			Type:             "健行",
			CreatedByUser:    &users[1],
			Participants: &EventsParticipants{
				LatestThreeUser: users,
				RemainNumber:    9,
			},
			CreatedAt: strPtr("2025-03-15T12:00:00Z"),
		},
	}

	return events, nil
}

func (c *Client) GetIdeas(ctx context.Context, ideaPast string, ideaPeriodBegin string, ideaPeriodEnd string) ([]Idea, error) {
	ideas := []Idea{
		{
			ID:          "1",
			Title:       "測試點子1",
			Description: "這是一個測試點子的描述",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "2",
			Title:       "測試點子2",
			Description: "這是另一個測試點子的描述",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			ID:          "3",
			Title:       "測試點子3",
			Description: "這是第三個測試點子的描述",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
	}
	return ideas, nil
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}
