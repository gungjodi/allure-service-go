package models

type ExecutorInfo struct {
	BuildName  string `json:"buildName"`
	BuildOrder int    `json:"buildOrder"`
	BuildURL   string `json:"buildUrl"`
	Name       string `json:"name"`
	ReportName string `json:"reportName"`
	ReportURL  string `json:"reportUrl"`
	Type       string `json:"type"`
}
