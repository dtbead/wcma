# contributing 
### needed tools
- [sqlc](https://sqlc.dev/) `./sqlc generate` (to generate sql queries)
- [mockgen](https://github.com/uber-go/mock) `./mockgen -source="storage/storage.go" > helper/testing/mock/storage/storage.go` (to generate mock packages)
- [postgres](http://postgresql.org/) ([dockerized](https://hub.docker.com/_/postgres/))

# building
run `go build` on `main.go`

# basic running
1. have a postgres instance ready
2. execute `internal\storage\postgres\schema.sql` in postgres (database migration, automatic scehma initialization, etc will be added later)
3. modify `wc_main_pg` string in `main.go` to your postgres instance
4. build and run

# notes
none of this is useful nor ready for anything meaningful whatsoever in its current state. 
> These things, they take time.