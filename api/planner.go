/*
	Copyright © 2022 Tom Lister tom@tomlister.net
*/
package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/tomlister/blik/config"
)

type CanvasPlannerItem struct {
	ContextType     string          `json:"context_type"`
	CourseID        int             `json:"course_id"`
	PlannableID     int             `json:"plannable_id"`
	PlannerOverride interface{}     `json:"planner_override"`
	PlannableType   string          `json:"plannable_type"`
	NewActivity     bool            `json:"new_activity"`
	Submissions     interface{}     `json:"submissions"`
	PlannableDate   time.Time       `json:"plannable_date"`
	Plannable       json.RawMessage `json:"plannable,omitempty"`
	HTMLURL         string          `json:"html_url"`
	ContextName     string          `json:"context_name"`
	ContextImage    string          `json:"context_image"`
}

type PlannableCalendarEvent struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	LocationName     string      `json:"location_name"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	AllDay           bool        `json:"all_day"`
	LocationAddress  interface{} `json:"location_address"`
	Description      string      `json:"description"`
	StartAt          time.Time   `json:"start_at"`
	EndAt            time.Time   `json:"end_at"`
	OnlineMeetingURL string      `json:"online_meeting_url"`
}

type PlannableAssignment struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	AssignmentID   int       `json:"assignment_id"`
	PointsPossible float64   `json:"points_possible"`
	DueAt          time.Time `json:"due_at"`
}

type PlannableQuiz struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PointsPossible float64   `json:"points_possible"`
	DueAt          time.Time `json:"due_at"`
}

func GetPlanner(cfg *config.Config) ([]CanvasPlannerItem, error) {
	var items []CanvasPlannerItem

	now := time.Now()

	u := url.URL{Scheme: "https", Host: cfg.Endpoint, Path: "/api/v1/planner/items"}

	q := u.Query()
	q.Add("start_date", now.Format("2006-01-02"))
	q.Add("per_page", "10")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return items, err
	}

	req.Header.Add("Authorization", "Bearer "+cfg.Key)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return items, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&items)
	if err != nil {
		return items, err
	}

	return items, nil
}
