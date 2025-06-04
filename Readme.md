# Restic Exporter

Restic Exporter is a tool for monitoring Restic repositories with Prometheus. Key features:
•	Collects statistics from Restic repositories (size, number of snapshots, etc.)
•	Exports metrics in Prometheus format
•	Manages repository passwords via CLI (in-progress)
•	Provides an HTTP API for accessing metrics



## Usage
### Run exporter
```
./restic-exporter serve
```
This command run http server on with /metrics endpoint.
Listened port set in .env file. 8085 by default.
This endpoint provide prometheus formated metrics of all restic repos that located in base path (Set it in .env file)


Set password from files for all repos
```
./restic-exporter set-password-cmd <directory> <new-password>
```
#### Arguments:
 - directory     Path to directory containing access files (*.access.*)
 - new-password  New password to set for all repositories


### Components
| Component       |                    Vendor                    |
|-----------------|:--------------------------------------------:|
| Config provider |   [Viper](https://github.com/spf13/viper)    |
| DI Container    |  [Sarulabs](https://github.com/sarulabs/di)  |
| Logger          |    [Zap](https://github.com/uber-go/zap)     |
| Mocks           | [Mockery](https://github.com/vektra/mockery) |
| CLI commands    |   [Cobra](https://github.com/spf13/cobra)    |
| Http routing    |   [Mux](github.com/gorilla/mux)    |



## Testing

### Test coverage

Run 
```
go test ./... -coverprofile cover.out
```
to collect coverage statistics.

Run
```
go tool cover -html=cover.out
```
to build ui statistic page

 
### Performance testing
You can run concurrent calls of grpc method by using [ghz utility](https://ghz.sh/docs/examples).   
```
ghz --insecure --proto internal/infrastructure/ports/grpc/proto/new.proto --call grpc.News.List -d '{"Author": {"Id": "44266dc6-18d0-46bd-a2b5-238de53db2cb"}, "Page": 1, "Query": "", "Sort": "ASC", "Status": "ACTIVE"}' -n 2000 -c 20 --connections=10 --debug ./debug.json   0.0.0.0:50051
```

### Mocks generation
```
mockery --all
```
Generate mock structures in pkg dir 
