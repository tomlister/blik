/*
	Copyright © 2022 Tom Lister tom@tomlister.net
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/tomlister/blik/api"
	"github.com/tomlister/blik/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blik",
	Short: "Overview of Canvas Tasks",
	Long:  "Overview of Canvas Tasks",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ReadConfig()
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("Please authenticate before running this command, refer to the help documentation using \"blik help\".")
				return
			}
			log.Fatal(err)
		}

		courses, err := api.GetCourses(cfg)
		if err != nil {
			log.Fatal(err)
		}

		coursesMapped := api.MapCourseToID(courses)

		items, err := api.GetPlanner(cfg)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			switch item.PlannableType {
			case "calendar_event":
				/*var event PlannableCalendarEvent
				err := json.Unmarshal(item.Plannable, &event)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("[EVENT] %s->%s %s\n", event.StartAt.Format(time.RFC1123), event.EndAt.Format(time.RFC1123), event.Title)*/
			case "assignment":
				var assignment api.PlannableAssignment
				err := json.Unmarshal(item.Plannable, &assignment)
				if err != nil {
					log.Fatal(err)
				}
				title := color.Style{color.FgCyan, color.OpBold}.Render(assignment.Title)
				url := url.URL{Scheme: "https", Host: cfg.Endpoint, Path: item.HTMLURL}
				fmt.Printf("[ASSIGNMENT] %s | %s\n> Due at: %s\n> Worth: %.0f\n> URL: %s\n\n", title, coursesMapped[item.CourseID].Name, assignment.DueAt.Format(time.RFC1123), assignment.PointsPossible, url.String())
			case "quiz":
				var quiz api.PlannableQuiz
				err := json.Unmarshal(item.Plannable, &quiz)
				if err != nil {
					log.Fatal(err)
				}
				title := color.Style{color.FgCyan, color.OpBold}.Render(quiz.Title)
				url := url.URL{Scheme: "https", Host: cfg.Endpoint, Path: item.HTMLURL}
				fmt.Printf("[QUIZ] %s | %s\n> Due at: %s\n> Worth: %.0f\n> URL: %s\n\n", title, coursesMapped[item.CourseID].Name, quiz.DueAt.Format(time.RFC1123), quiz.PointsPossible, url.String())
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
