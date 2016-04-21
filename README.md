Site-Checker Site Availability Checker

Provides really basic site availability checks, which run in parallel so you can check multiple sites at once.
Configure expected response code + page content (we'd expect to see a 200 with "webok dbok" if everything is good for example).

Configuration via JSON file, place in same directory as executable.
Will also create logs & state directories - log to file, access log style, plus basic email alerts if enabled.
Recommended to have executable & conf.json in /opt/site-checker

Runs configured checks then exists - we run it as a cron, twice per minute.
