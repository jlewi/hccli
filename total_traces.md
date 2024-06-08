# Analyze traffic

How many users did we have in the last 7 days

```bash
/Users/jlewi/git_hccli/hccli nltoq --nlq=\"Count the number of traces for the last 7 days broken down by regions\" --dataset=foyle
```

```bash
/Users/jlewi/git_hccli/hccli createquery --query='{\"breakdowns\": [\"http.method\"], \"calculations\": [{\"op\": \"COUNT\"}], \"filters\": [{\"column\": \"http.method\", \"op\": \"exists\"}], \"orders\": [{\"op\": \"COUNT\", \"order\": \"descending\"}], \"time_range\": 604800}' --dataset=foyle
```
