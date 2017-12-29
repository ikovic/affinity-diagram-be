package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jiraInstance := os.Getenv("JIRA_INSTANCE")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraPassword := os.Getenv("JIRA_PASSWORD")

	jiraClient, err := jira.NewClient(nil, jiraInstance)
	if err != nil {
		panic(err)
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(jiraUsername, jiraPassword)
	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		panic(err)
	}

	projectList, _, _ := jiraClient.Project.GetList()

	for _, project := range *projectList {
		fmt.Printf("Result: %v\n", project.Name)
	}

	router := httprouter.New()
	router.GET("/", index)

	log.Fatal(http.ListenAndServe(":8080", router))
}
