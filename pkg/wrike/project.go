package wrike

import (
	"fmt"
	"time"
)

// Projects contains the kind of the project and the data related to it.
type Projects struct {
	Kind string    `json:"kind"`
	Data []Project `json:"data"`
}

// Project holds the details about a single project.
type Project struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Children    []string `json:"children"`
	ChildIds    []string `json:"childIds"`
	Scope       string   `json:"scope"`
	Project     struct {
		AuthorID       string    `json:"authorId"`
		OwnerIds       []string  `json:"ownerIds"`
		CustomStatusID string    `json:"customStatusId"`
		StartDate      string    `json:"startDate"`
		EndDate        string    `json:"endDate"`
		CreatedDate    time.Time `json:"createdDate"`
	} `json:"project,omitempty"`
}

// AllProjectMap is a map with project id as key and Project struct as value
type AllProjectMap map[string]Project

// GetProjectsByIds gets projects for given project IDs
func (pm *AllProjectMap) GetProjectsByIds(projectIDs []string) []Project {
	var projects []Project
	for _, id := range projectIDs {
		projects = append(projects, (*pm)[id])
	}
	return projects
}

// GetAllProjects fetches all projects and returns them as a map where projectId is the key
func (w *Client) GetAllProjects() (AllProjectMap, error) {
	urlQuery := map[string]string{
		"deleted": "false",
		"project": "true",
		"fields":  `["description"]`,
	}

	var projects Projects

	resp, err := w.httpClient.R().SetQueryParams(urlQuery).
		SetResult(&projects).
		Get("/spaces/" + w.spaceId + "/folders")

	if err != nil {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), resp.String())
	}

	projectMap := AllProjectMap{}
	for _, project := range projects.Data {
		projectMap[project.ID] = project
	}

	return projectMap, nil
}

// GetProjectsByLink fetches projects & folders filtered by link
func (w *Client) GetProjectsByLink(link string, urlQuery map[string]string) (Projects, error) {
	if urlQuery == nil {
		urlQuery = map[string]string{
			"deleted": "false",
		}
	}
	if len(link) > 0 {
		urlQuery["permalink"] = link
	}

	var projects Projects
	resp, err := w.httpClient.R().SetQueryParams(urlQuery).
		SetResult(&projects).
		Get("/folders")

	if err != nil {
		return projects, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), resp.String())
	}

	return projects, nil
}
