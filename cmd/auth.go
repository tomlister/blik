/*
	Copyright © 2022 Tom Lister tom@tomlister.net
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomlister/blik/config"
)

type CanvasUser struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	SortableName    string    `json:"sortable_name"`
	ShortName       string    `json:"short_name"`
	AvatarURL       string    `json:"avatar_url"`
	EffectiveLocale string    `json:"effective_locale"`
}

func testKey(endpoint string, key string) (string, error) {
	u := url.URL{Scheme: "https", Host: endpoint, Path: "/api/v1/users/self"}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+key)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var user CanvasUser
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.Name, nil
}

func promptKey() (string, error) {
	fmt.Print("Enter Key: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimSuffix(line, "\n")), nil
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth [endpoint] [key]",
	Short: "Authenticates your canvas account",
	Long:  "Authenticates your canvas account",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]
		key, err := promptKey()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Testing key...")
		name, err := testKey(endpoint, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Welcome, %s.\n", name)
		cfg := config.NewConfig(endpoint, key)
		err = cfg.Save()
		if err != nil {
			log.Fatal("An error occurred while saving the configuration.")
		}
		fmt.Println("Configuration successfully saved.")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
