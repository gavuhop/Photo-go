# Initialize variables
PACKAGE_PATH="./..."
has_f_flag=false
COVERAGE_HTML="coverage.html"
# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            echo "Usage: $0 [-p <package_path>] [-f]"
            echo "   Example: $0 -p cmd/api/main_test.go -f"
            echo "       \"-p\" specifies the package path (default: ./...)"
            echo "       \"-f\" will run all tests and show full output"
            echo "       without \"-f\" will run all tests and show coverage report"
            exit 0
            ;;
        -p)
            PACKAGE_PATH="./"$2"/..."
            shift 2
            ;;
        -f)
            has_f_flag=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done
BASE_PATH=$(pwd)
echo "BASE_PATH: $BASE_PATH"
export BASE_PATH=$BASE_PATH

eval $(go run config/env/secret_base64.go -p $BASE_PATH/secret_file_dev.json)

export NOT_DEBUG=true
export ENVIRONMENT=test

go clean -testcache

if [ "$has_f_flag" = true ]; then
    go test -v -coverpkg=$PACKAGE_PATH -coverprofile=profile.cov $PACKAGE_PATH
else
    go test -v -coverpkg=$PACKAGE_PATH -coverprofile=profile.cov $PACKAGE_PATH 2>&1 | grep -A 5 -E "FAIL|--- FAIL"
fi
go tool cover -func profile.cov
go tool cover -html=profile.cov -o $COVERAGE_HTML
