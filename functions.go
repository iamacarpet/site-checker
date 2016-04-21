package main

import (
    "os"
    "path/filepath"
    "crypto/md5"
    "encoding/hex"
)

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func basePath() (string) {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        return ""
    }

    return dir + string(os.PathSeparator)
}

func md5sum(text string) (string) {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}
