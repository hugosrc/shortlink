package domain

import "time"

type LinkMetrics struct {
	ShortURL       string    `json:"short_url"`
	OriginalURL    string    `json:"original_url"`
	IPAddress      string    `json:"ip_address"`
	Referer        string    `json:"referer"`
	Device         string    `json:"device"`
	OS             string    `json:"os"`
	OSVersion      string    `json:"os_version"`
	UserAgent      string    `json:"user_agent"`
	UserAgentName  string    `json:"user_agent_name"`
	Version        string    `json:"version"`
	AcceptLanguage string    `json:"accept_language"`
	AccessTime     time.Time `json:"access_time"`
}
