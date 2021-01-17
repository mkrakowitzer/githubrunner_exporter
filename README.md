# A prometheus exporter for github runner metrics

[![Go Report Card](https://goreportcard.com/badge/github.com/mkrakowitzer/githubrunner_exporter)](https://goreportcard.com/report/github.com/mkrakowitzer/githubrunner_exporter)

# Parameters

You can set the configuration options at runtime, as environment variables or yaml configuration file.

#### Environment variables

| Param               | Description
| ------------------- | ---------------------------------- |
| `GITHUB_INTERVAL`   | Interval to query the github API   |
| `GITHUB_ORG`        | Organisation name                  |
| `GITHUB_TOKEN`      | GitHub Personal Access Token       |

#### Runtime
To see all available configuration flags:
```console
./githubrunner_exporter --interval 15 --token XXX_PERSONAL_ACCESS_TOKEN_XXX --org myorgname
```

#### YAML configuration file
$HOME/.githubrunner_exporter.yaml
```yaml
---
interval: 15
org: camelotls
token: XXX_PERSONAL_ACCESS_TOKEN_XXX
```

The exporter makes use of etags when quering the API. When data has not changed GitHub returns a 304 response does not count against your Rate Limit See https://docs.github.com/en/rest/overview/resources-in-the-rest-api#conditional-requests.

# Metrics

| Name                       | Description                                                   |
| -------------------------- | ------------------------------------------------------------- |
| github_runner_busy         | The runner has an active job running.                         |
| github_runner_status       | The runner is online/offline.                                 |
| github_ratelimit_limit     | The total number of allowed API calls for your user/token     |
| github_ratelimit_remaining | Remaining allowed API calls                                   |
| github_ratelimit_used      | Total number of calls made                                    |
| github_ratelimit           | When Rate limit will reset (Every hour). Epoch Time           |

#### github_runner_status

| Value | Description |
| ------| ----------- |
| 1     | online      |
| 0     | offline     |

#### github_runner_busy

| Value | Description |
| ------| ----------- |
| 1     | busy        |
| 0     | offline     |

```
# HELP github_ratelimit Time until rate limit resets (epoch)
# TYPE github_ratelimit gauge
github_ratelimit 1.610879916e+09
# HELP github_ratelimit_limit Total number of calls allowed
# TYPE github_ratelimit_limit gauge
github_ratelimit_limit 5000
# HELP github_ratelimit_remaining Number of calls remaining
# TYPE github_ratelimit_remaining gauge
github_ratelimit_remaining 4994
# HELP github_ratelimit_used Number of used calls of your rate limit
# TYPE github_ratelimit_used gauge
github_ratelimit_used 6
# HELP github_runner_busy runner busy
# TYPE github_runner_busy gauge
github_runner_busy{id="888",name="aws-runner-01",os="linux"} 1
github_runner_busy{id="889",name="aws-runner-02",os="linux"} 1
# HELP github_runner_status runner status
# TYPE github_runner_status gauge
github_runner_status{id="888",name="aws-runner-01",os="linux"} 1
github_runner_status{id="889",name="aws-runner-02",os="linux"} 1
```

# Installation

## Docker

You can run the container with the following commands
```console
docker login ghcr.io
docker pull ghcr.io/mkrakowitzer/githubrunner_exporter
docker run --env GITHUB_INTERVAL=15 --env GITHUB_ORG=orgname --env GITHUB_TOKEN=XXX_TOKEN_XXX \
    -d -p 9090 --name deleteme ghcr.io/mkrakowitzer/githubrunner_exporter
```

# Alert rules

## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)

Building:

```console
git clone git@github.com:mkrakowitzer/githubrunner_exporter.git
cd githubrunner_exporter
go build -o githubrunner_exporter
./githubrunner_exporter <flags>
```

To see all available configuration flags:
```console
./githubrunner_exporter -h
```

## Running tests

```console
make test
```
