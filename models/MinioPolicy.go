package models

type MinioPolicy struct {
	Version   string `json:"Version"`
	Statement []struct {
		Effect   string   `json:"Effect"`
		Action   []string `json:"Action"`
		Resource []string `json:"Resource"`
	} `json:"Statement"`
}
