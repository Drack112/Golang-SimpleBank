case "$1" in
    -s | --short)
        case "$2" in
            -c | --coverage) echo "Run only unit tests (with coverage)"
                go test -coverpkg ./... ./... -coverprofile=coverage.txt
                goveralls -coverprofile=coverage.txt -service semaphore -repotoken $COVERALLS_TOKEN
            ;;
            *) echo "Run only unit tests"
                go test -v -short ./...
            ;;
        esac
    ;;
    -i | --integration)  echo  "Run only integration tests"
        go test -v -run Integration ./...
    ;;
    *) echo "Run all tests (with coverage)"
        go test -coverpkg ./... ./... -coverprofile=coverage.txt
        goveralls -coverprofile=coverage.txt -service semaphore -repotoken $COVERALLS_TOKEN
    ;;
esac
