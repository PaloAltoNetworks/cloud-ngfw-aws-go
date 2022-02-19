package rulestack

import (
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/permissions"
)

// List returns a list of objects.
func List(client api.Client, input ListInput) (ListOutput, error) {
	client.Log(http.MethodGet, "list rulestacks")

	var ans ListOutput
	_, err := client.Communicate(
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func Create(client api.Client, input Info) error {
	client.Log(http.MethodPost, "create rulestack: %s", input.Name)

	_, err := client.Communicate(
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func Read(client api.Client, input ReadInput) (ReadOutput, error) {
	name := input.Name
	input.Name = ""

	client.Log(http.MethodGet, "describe rulestack: %s", name)

	var ans ReadOutput
	_, err := client.Communicate(
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func Update(client api.Client, input Info) error {
	name := input.Name
	input.Name = ""

	client.Log(http.MethodPut, "updating rulestack: %s", name)

	_, err := client.Communicate(
		permissions.Rulestack,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func Delete(client api.Client, name string) error {
	client.Log(http.MethodDelete, "delete rulestack: %s", name)

	_, err := client.Communicate(
		permissions.Rulestack,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		nil,
		nil,
	)

	return err
}
