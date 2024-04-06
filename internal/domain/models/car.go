package models

type Car struct {
    RegNum string `json:"regNum"`
    Mark   string `json:"mark"`
    Model  string `json:"model"`
    Year   int64  `json:"year,omitempty"`
    Owner  People `json:"owner"`
}