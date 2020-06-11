***********************************

# HTTP log monitoring console program

***********************************

# Use Case
Consume an actively written-to CLF HTTP access log (https://en.wikipedia.org/wiki/Common_Log_Format). 
- It should default to reading /tmp/access.log and be overrideable
Example log lines:

```127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123

127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 200 234

127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 200 34

127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12
```

- Display stats every 10s about the traffic during those 10s: the sections of the web site with the most hits, as well as interesting summary statistics on the traffic as a whole. A section is defined as being what's before the second '/' in the resource section of the log line. For example, the section for "/pages/create" is "/pages"
- Make sure a user can keep the app running and monitor the log file continuously
- Whenever total traffic for the past 2 minutes exceeds a certain number on average, print or display a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”. The default threshold should be 10 requests per second, and should be overridable
- Whenever the total traffic drops again below that value on average for the past 2 minutes, print or display another message detailing when the alert recovered
- Write a test for the alerting logic
- Explain how you’d improve on this application design

# Usage

Everything is in the monitor package.
A NewLauncher() func is provided to initialize a launcher
```
l := NewLauncher()
```

It is possible to specify:
- filepath
- alertign treshold
on the command line.

There is also an functional option to fully configure the launcher:
```
l := monitor.NewLauncher(monitor.DefaultConfig(filePath, treshold))
```
The ***DefaultConfig(filePath, treshold)*** is a helper to configure the launcher, but really what's needed is an Options:
```
type Options func(*Launcher)
```

The iuse the launch funtion
```
l.Launch(io.Writer)
```
You may specify os.Stdout as default.

# Output

Every reporting interval (default to 10s) a report is printed to the target io.Writer:
- pages hits rank (5 most hit pages)
- protocol rank (5 most used protocols)
- Status code rank (5 most)

Every alerting interval (default to 120s) an alert is printed to the target io.Writer 
if traffic on average for the interval is above treshold(default to 10 per Second).
Another alert is sent when traffic recover.

# Improvement
- Did not have enough time for the printing function. Would be nice to:
  - make it more friendly and shiny
  - transfer it to its own structure (better for maintanaibility)
  - use a concurrent function to print... We possibly could have lots of data to process
- Provide a Launcher interface so oother developer can provide their own version
- Add a statistic about response size( the reponse size is already kept in the entry struct)
- refactor session struct to make it more configurable
