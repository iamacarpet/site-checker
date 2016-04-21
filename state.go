package main

import (
    "os"
    "encoding/json"
)

func loadState(hash string) (*StatusStorage, error) {
    file, _ := os.Open(basePath() + "state" + string(os.PathSeparator) + hash + ".json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    status := &StatusStorage{}
    err := decoder.Decode(status)
    return status, err
}

func saveState(hash string, site *SiteConfiguration, resp *SiteResponse) {
    statePath := basePath() + "state"
    if exist, _ := exists(statePath); !exist {
        os.Mkdir(statePath, 0755)
    }
    file, _ := os.Create(statePath + string(os.PathSeparator) + hash + ".json")
    defer file.Close()
    encoder := json.NewEncoder(file)

    state := &StatusStorage{}
    if site != nil {
        state.Site = *site
    }
    if resp != nil {
        state.Response = *resp
    }

    _ = encoder.Encode(state)
}
