#!/bin/sh

# Run the tests
go test ./tests/...

# Check if tests passed
if [ $? -eq 0 ]; then
    echo "All tests passed. Starting the application..."
    /apigwservice
else
    echo "Tests failed. Exiting."
    exit 1
fi

