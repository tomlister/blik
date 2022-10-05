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

type CanvasCourse struct {
	ID                          int         `json:"id"`
	Name                        string      `json:"name"`
	AccountID                   int         `json:"account_id"`
	UUID                        string      `json:"uuid"`
	StartAt                     time.Time   `json:"start_at"`
	GradingStandardID           int         `json:"grading_standard_id"`
	IsPublic                    bool        `json:"is_public"`
	CreatedAt                   time.Time   `json:"created_at"`
	CourseCode                  string      `json:"course_code"`
	DefaultView                 string      `json:"default_view"`
	RootAccountID               int         `json:"root_account_id"`
	EnrollmentTermID            int         `json:"enrollment_term_id"`
	License                     string      `json:"license"`
	GradePassbackSetting        interface{} `json:"grade_passback_setting"`
	EndAt                       interface{} `json:"end_at"`
	PublicSyllabus              bool        `json:"public_syllabus"`
	PublicSyllabusToAuth        bool        `json:"public_syllabus_to_auth"`
	StorageQuotaMb              int         `json:"storage_quota_mb"`
	IsPublicToAuthUsers         bool        `json:"is_public_to_auth_users"`
	HomeroomCourse              bool        `json:"homeroom_course"`
	CourseColor                 interface{} `json:"course_color"`
	FriendlyName                interface{} `json:"friendly_name"`
	ApplyAssignmentGroupWeights bool        `json:"apply_assignment_group_weights"`
	Calendar                    struct {
		Ics string `json:"ics"`
	} `json:"calendar"`
	TimeZone    string `json:"time_zone"`
	Blueprint   bool   `json:"blueprint"`
	Template    bool   `json:"template"`
	Enrollments []struct {
		Type                           string `json:"type"`
		Role                           string `json:"role"`
		RoleID                         int    `json:"role_id"`
		UserID                         int    `json:"user_id"`
		EnrollmentState                string `json:"enrollment_state"`
		LimitPrivilegesToCourseSection bool   `json:"limit_privileges_to_course_section"`
	} `json:"enrollments"`
	HideFinalGrades                  bool   `json:"hide_final_grades"`
	WorkflowState                    string `json:"workflow_state"`
	RestrictEnrollmentsToCourseDates bool   `json:"restrict_enrollments_to_course_dates"`
	OverriddenCourseVisibility       string `json:"overridden_course_visibility,omitempty"`
}

func GetCourses(cfg *config.Config) ([]CanvasCourse, error) {
	var courses []CanvasCourse
	u := url.URL{Scheme: "https", Host: cfg.Endpoint, Path: "/api/v1/courses"}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return courses, err
	}

	req.Header.Add("Authorization", "Bearer "+cfg.Key)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return courses, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&courses)
	if err != nil {
		return courses, err
	}

	return courses, nil
}

func MapCourseToID(courses []CanvasCourse) map[int]CanvasCourse {
	mapped := map[int]CanvasCourse{}
	for _, c := range courses {
		mapped[c.ID] = c
	}
	return mapped
}
