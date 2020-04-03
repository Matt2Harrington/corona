<a href="http://fvcproductions.com"><img src="https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/240/apple/237/microbe_1f9a0.png" title="FVCproductions" alt="FVCproductions"></a>

<!-- [![FVCproductions](https://avatars1.githubusercontent.com/u/4284691?v=3&s=200)](http://fvcproductions.com) -->

COVID-19

# Covid-19 (Coronavirus Data)

><a href="github.com/NovelCOVID/API">Using NovelCOVID API</a>

> Taking data and inserting into Postgres Database for future use.

## Prerequisites

- Download and install <a href="https://www.postgresql.org/download/">Postgres</a>
- Use data application (such as DataGrip) to connect to Postgres locally
- Create a config.yaml as follows and fill in values:

```yaml
    host: <host>
    username: <username>
    port: <port>
    databaseName: <name>
```
- Run provided scripts to create tables
- Run insertion scripts to load existing data

## Running
- Have Postgres running locally
- Navigate to directory of repo and run `go run corona.go` to pull new data
- (Data usually gets updated 25-30 minutes on API call)

## Storage
- Location data is stored in the `info` table
- `data_id` links to the id in the data table for specific data records
