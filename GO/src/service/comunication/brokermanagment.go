package comunication

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BrokerManagement struct {
	rabbitMQURL   string
	adminUsername string
	adminPassword string
	auth          string
	client        *http.Client
}

func (m *BrokerManagement) CreateUser(username, password string) error {
	createUserURL := fmt.Sprintf("%s/api/users/%s", m.rabbitMQURL, username)
	userJSON := []byte(fmt.Sprintf(`{"password":"%s", "tags":"management"}`, password))
	return m.sendRequest("PUT", createUserURL, userJSON)
}

func (m *BrokerManagement) SetPermission(username string) error {
	setPermissionsURL := m.rabbitMQURL + "/api/permissions/%2F/" + username
	permissionsJSON := []byte(`{"configure":".*","write":".*","read":".*"}`)
	return m.sendRequest("PUT", setPermissionsURL, permissionsJSON)
}

func (m *BrokerManagement) SetTopicPermission(username string, topics []string) error {
	regex := "null"

	if len(topics) > 0 {
		// Customize patterns based on your specific requirements
		regex = strings.Join(topics, "|")
	}
	// Set permissions for the user
	setPermissionsURL := m.rabbitMQURL + "/api/topic-permissions/%2F/" + username
	permissionsJSON := []byte(fmt.Sprintf(`{"exchange":"amq.topic","write":"%s","read":"%s"}`, regex, regex))
	return m.sendRequest("PUT", setPermissionsURL, permissionsJSON)
}
func NewBrokerManagement(rabbitMQURL, adminUsername, adminPassword string) *BrokerManagement {
	// Encode adminUsername and adminPassword for Basic Authentication
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", adminUsername, adminPassword)))
	return &BrokerManagement{
		rabbitMQURL: rabbitMQURL, adminUsername: adminUsername, adminPassword: adminPassword, auth: auth,
		client: &http.Client{},
	}
}

func (m *BrokerManagement) sendRequest(method, url string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+m.auth)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("Response Body: %s\n", responseBody)

	// Handle the response body or any additional processing here

	return nil
}
