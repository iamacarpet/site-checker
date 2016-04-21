package main

import (
    "os"
    "fmt"
    "time"
    "bytes"
    "syscall"
    "net/smtp"
    "encoding/json"
)

type SiteResponse struct {
    Code        int                 `json:"code"`
    Body        string              `json:"body"`
    OK          bool                `json:"OK"`
}

type StatusStorage struct {
    Site        SiteConfiguration   `json:"site"`
    Response    SiteResponse        `json:"response"`
}

func logSiteLine(site *SiteConfiguration, msg string, resp *SiteResponse, dur time.Duration){
    // Generate our site hash & log file name...
    hashTemp, _ := json.Marshal(site)
    siteHash := md5sum(string(hashTemp))
    // First off, load our status file, check if we're the same as previous.
    prevState, _ := loadState(siteHash)

    if resp == nil {
        resp = &SiteResponse{ OK: false }
    }
    if ! resp.OK {
        if prevState.Response.OK != false {
            msg += " - Email Sent"
        }
    } else {
        if prevState.Response.OK != true {
            msg += " - Email Sent"
        }
    }

    logLine(fmt.Sprintf("%03d %10.6f %-50s %s", resp.Code, dur.Seconds(), site.URL, msg))
    saveState(siteHash, site, resp)

    if ! resp.OK {
        if prevState.Response.OK != false {
            sendEmail(false, msg, site, resp)
        }
    } else {
        if prevState.Response.OK != true {
            sendEmail(true, msg, site, resp)
        }
    }
}

func logLine(data string) {
    t := time.Now()
    fileFolder := basePath() + "logs"
    var fileName string = fileFolder + string(os.PathSeparator) + "site-checker-" + t.Format("2006-01-02") + ".log"

    err := os.MkdirAll(fileFolder, 0755)
    if err != nil {
        return
    }
    f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0655)
    if err != nil {
        return
    }
    defer f.Close()

    syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
    defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

    f.WriteString(fmt.Sprintf("%s: %s", t.Format(time.RFC1123), data) + "\r\n")
}

func sendEmail(recovery bool, msg string, site *SiteConfiguration, resp *SiteResponse){
    if ! config.EmailConfig.SendEmails {
        return
    }

    for _, email := range config.EmailConfig.AlertsTo {
        emailBody := "From: " + hostname + " <" + config.EmailConfig.FromEmail + ">\r\n"
        emailBody += "To: Alert User <" + email + ">\r\n"
        emailBody += "Subject: "
        if recovery {
            emailBody += "RECOVERY - "
        } else {
            emailBody += "PROBLEM - "
        }
        emailBody += site.FriendlyName + "\r\n"
        emailBody += "Content-Type: text/plain; charset=UTF-8\r\n"
        emailBody += "\r\n"

        emailBody += "Message: " + msg + "\r\n"
        if resp != nil {
            emailBody += "HTTP Code: " + fmt.Sprintf("%03d", resp.Code) + "\r\n"
            emailBody += "HTTP Body: " + resp.Body + "\r\n"
        }

        func (){
            c, err := smtp.Dial(config.EmailConfig.SMTPServer + ":25")
            if err != nil {
                return
            }
            defer c.Close()

            c.Mail(config.EmailConfig.FromEmail)
            c.Rcpt(email)

            wc, err := c.Data()
            if err != nil {
                return
            }
            defer wc.Close()

            buf := bytes.NewBufferString(emailBody)
            if _, err = buf.WriteTo(wc); err != nil {
                return
            }
        }()
    }
}
