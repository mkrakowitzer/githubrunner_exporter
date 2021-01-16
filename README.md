# A prometheus exporter for github runner metrics

# Parameters
`GITHUB_INTERVAL` - Interval to query the github API
`GITHUB_ORG` - Organisation name
`GITHUB_TOKEN` - GitHub Token

The exporter makes use of etags when quering the API. When data has not changed GitHub returns a 304 response does not count against your Rate Limit See https://docs.github.com/en/rest/overview/resources-in-the-rest-api#conditional-requests.

