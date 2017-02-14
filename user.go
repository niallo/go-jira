package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// UserService handles users for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user
type UserService struct {
	client *Client
}

// User represents a JIRA user.
type User struct {
	Self            string     `json:"self,omitempty" structs:"self,omitempty"`
	Name            string     `json:"name,omitempty" structs:"name,omitempty"`
	Password        string     `json:"-"`
	Key             string     `json:"key,omitempty" structs:"key,omitempty"`
	EmailAddress    string     `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
	AvatarUrls      AvatarUrls `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	DisplayName     string     `json:"displayName,omitempty" structs:"displayName,omitempty"`
	Active          bool       `json:"active,omitempty" structs:"active,omitempty"`
	TimeZone        string     `json:"timeZone,omitempty" structs:"timeZone,omitempty"`
	ApplicationKeys []string   `json:"applicationKeys,omitempty" structs:"applicationKeys,omitempty"`
}

// Get gets user info from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
func (s *UserService) Get(username string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?username=%s", username)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Create creates an user in JIRA.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-createUser
func (s *UserService) Create(user *User) (*User, *Response, error) {
	apiEndpoint := "/rest/api/2/user"
	req, err := s.client.NewRequest("POST", apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseUser := new(User)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	err = json.Unmarshal(data, responseUser)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return responseUser, resp, nil
}

// SearchOptions specifies the optional parameters to various List methods that
// support pagination.
// Pagination is used for the JIRA REST APIs to conserve server resources and limit
// response size for resources that return potentially large collection of items.
// A request to a pages API will result in a values array wrapped in a JSON object with some paging metadata
// Default Pagination options
type FindUsersOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults int
	// IncludeActive: If true, then active users are included in the results. Default: true.
	IncludeActive bool
	// IncludeInactive: If true, then inactive users are included in the results. Default: false.
	IncludeInactive bool
	// Property: A query string used to search by property.
	// Property key cannot contain dot or equal sign, value cannot be JSONObject.
	// Example: for following property value: {"something":{"nested":1,"other":2}},
	// you can search: propertyKey.something.nested=1.
	Property string
}

// Search will search for users according to the username and options.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-findUsers
func (s *UserService) FindUsers(username string, options *FindUsersOptions) ([]User, *Response, error) {
	var u string
	if options == nil {
		u = fmt.Sprintf("rest/api/2/user/search?username=%s", username)
	} else {
		u = fmt.Sprintf(
			"rest/api/2/user/search?username=%s&startAt=%d&maxResults=%d"+
				"&includeActive=%t&includeInactive=%t&Property=%s",
			url.QueryEscape(username), options.StartAt, options.MaxResults,
			options.IncludeActive, options.IncludeInactive,
			url.QueryEscape(options.Property))
	}

	users := []User{}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return []User{}, nil, err
	}

	resp, err := s.client.Do(req, &users)
	return users, resp, err
}
