package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gojira "github.com/andygrunwald/go-jira"
	jira "github.com/ikovic/affinity-diagram-be/jira"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type SimpleIssue struct {
	URL         string `json:"url" structs:"url"`
	Key         string `json:"key" structs:"key"`
	ID          string `json:"id" structs:"id"`
	Summary     string `json:"summary" structs:"summary"`
	Description string `json:"description" structs:"description"`
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getBoards(jiraClient *gojira.Client) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		options := gojira.BoardListOptions{BoardType: "scrum"}
		boardList, _, _ := jiraClient.Board.GetAllBoards(&options)
		bytes, _ := json.Marshal(boardList.Values)
		w.Write(bytes)
	}
}

func getBacklogIssues(jiraClient *gojira.Client) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		boardID := params.ByName("boardId")
		query := fmt.Sprintf("/rest/agile/1.0/board/%s/backlog", boardID)
		req, _ := jiraClient.NewRequest("GET", query, nil)

		result := new(gojira.IssuesInSprintResult)
		_, err := jiraClient.Do(req, result)
		if err != nil {
			panic(err)
		}

		simpleIssues := []SimpleIssue{}
		for _, issue := range result.Issues {
			simpleIssue := SimpleIssue{URL: issue.Self, Key: issue.Key, ID: issue.ID,
				Summary: issue.Fields.Summary, Description: issue.Fields.Description}
			simpleIssues = append(simpleIssues, simpleIssue)
		}

		bytes, _ := json.Marshal(simpleIssues)
		w.Write(bytes)
	}
}

func main() {
	loadEnv()

	jiraInstance := os.Getenv("JIRA_INSTANCE")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraPassword := os.Getenv("JIRA_PASSWORD")

	jiraClient := jira.GetClient(jiraInstance, jiraUsername, jiraPassword)

	router := httprouter.New()
	handler := cors.Default().Handler(router)
	router.GET("/", index)
	router.GET("/boards", getBoards(jiraClient))
	router.GET("/backlog/:boardId", getBacklogIssues(jiraClient))

	fmt.Println("Server listening at 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
