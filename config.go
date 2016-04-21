package main

import (
    "os"
    "encoding/json"
)

var config Configuration
var hostname string

type Configuration struct {
    EmailConfig     EmailConfiguration      `json:"emailConfig"`
    UserAgent       string                  `json:"userAgent"`
    Alias           string                  `json:"alias"`
    CheckSites      []SiteConfiguration     `json:"checkSites"`
    ProxyURL        string                  `json:"proxyURL"`
    Timeout         int                     `json:"timeout"`
}

type EmailConfiguration struct {
    FromEmail       string                  `json:"fromEmail"`
    AlertsTo        []string                `json:"alertsTo"`
    SMTPServer      string                  `json:"smtpServer"`
    SendEmails      bool                    `json:"enabled"`
}

type SiteConfiguration struct {
    UseProxy        bool                    `json:"useProxy"`
    URL             string                  `json:"url"`
    FriendlyName    string                  `json:"friendlyName"`
    ExpectedResp    string                  `json:"expectedContent"`
    ExpectedCode    int                     `json:"expectedCode"`
}

func loadConfig(){
    file, _ := os.Open(basePath() + "conf.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    config = Configuration{}
    err := decoder.Decode(&config)
    if err != nil {
        //saveConfig()
    }
}

func saveConfig(){
    file, _ := os.Create(basePath() + "conf.json")
    defer file.Close()
    encoder := json.NewEncoder(file)
    err := encoder.Encode(&config)
    if err != nil {
        panic(err)
    }
}
