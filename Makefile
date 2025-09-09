test:
	go test -v -coverprofile ./coverage.out -coverpkg ./...  ./...

view_cov:
	go tool cover -html=./coverage.out

view_cov_summary:
	go tool cover -func-./coverage.out