package wrike

import (
	"strings"
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
func (w *Client) GetAllProjects() AllProjectMap {
	urlQuery := map[string]string{
		"deleted": "false",
		"project": "true",
		"fields":  `["description"]`,
	}

	var projects Projects
	w.callAPI("/spaces/"+w.spaceId+"/folders", urlQuery, &projects)

	projectMap := AllProjectMap{}
	for _, project := range projects.Data {
		projectMap[project.ID] = project
	}

	return projectMap
}

// GetProjectsByLink fetches projects & folders filtered by link
func (w *Client) GetProjectsByLink(link string, urlQuery map[string]string) Projects {
	if urlQuery == nil {
		urlQuery = map[string]string{
			"deleted": "false",
		}
	}
	if len(link) > 0 {
		urlQuery["permalink"] = link
	}

	var projects Projects
	w.callAPI("/folders", urlQuery, &projects)

	return projects
}

// GetProjectsByIds fetches projects & folders filtered by Ids
func (w *Client) GetProjectsByIds(ids []string) Projects {
	var projects Projects
	if len(ids) > 0 {
		w.callAPI("/folders/"+strings.Join(ids, ","), nil, &projects)
	}

	return projects
}
