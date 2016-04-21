package main

import (
    "os"
    "fmt"
    "sync"
    "time"
    "strings"
    "net/url"
    "net/http"
    "io/ioutil"
)

func main(){
    hostname, _ = os.Hostname()
    loadConfig()

    var wg sync.WaitGroup

    wg.Add(len(config.CheckSites))

    for _, site := range config.CheckSites {
        go func(s SiteConfiguration){
            defer wg.Done()

            performCheck(&s)
        }(site)
    }

    wg.Wait()
}

func performCheck(site *SiteConfiguration){
    if site.ExpectedCode == 0 {
        site.ExpectedCode = 200
    }

    client := &http.Client{
        Timeout: time.Duration(config.Timeout) * time.Second,
    }
    if site.UseProxy && config.ProxyURL != "" {
        proxyURL, _ := url.Parse(config.ProxyURL)
        client.Transport = &http.Transport{ Proxy: http.ProxyURL(proxyURL) }
    }

    req, err := http.NewRequest("GET", site.URL, nil)
    if err != nil {
        logSiteLine(site, "INIT REQ FAILED", nil, time.Duration(0))
        return
    }
    req.Header.Set("User-Agent", strings.Replace(config.UserAgent, "%HOSTNAME%", hostname, -1))

    startTime := time.Now()
    resp, err := client.Do(req)
    timeTook  := time.Since(startTime)

    if err != nil {
        logSiteLine(site, fmt.Sprintf("REQ FAILED - ERR: %s", err.Error()), nil, timeTook)
        return
    }

    defer resp.Body.Close()
    siteResp := SiteResponse{}
    siteResp.Code = resp.StatusCode
    body, _ := ioutil.ReadAll(resp.Body)
    siteResp.Body = strings.TrimSpace(string(body))
    siteResp.OK = false

    if siteResp.Code != site.ExpectedCode {
        logSiteLine(site, "WRONG STATUS", &siteResp, timeTook)
        return
    }
    if site.ExpectedResp != "" {
        if siteResp.Body != site.ExpectedResp {
            logSiteLine(site, "WRONG BODY", &siteResp, timeTook)
            return
        }
    }

    siteResp.OK = true
    logSiteLine(site, "OK", &siteResp, timeTook)
    return
}
