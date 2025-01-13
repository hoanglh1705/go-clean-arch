#!/usr/bin/env bash

set -ex

# Generate swagger.json file
swagger generate spec --scan-models -w ./cmd/server/ -o ./swaggerui/swagger.json

if [[ "$AWS_LAMBDA_FUNCTION_NAME" != "" && "$STAGE" != "" && "$STAGE" != "development" ]]; then
    HOST=$(aws ssm get-parameters --name "/$AWS_LAMBDA_FUNCTION_NAME/$STAGE/host" --with-decryption | jq -r '.Parameters[0].Value')
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    # Replace HOST by corresponding env var
    sed -i '' -e "s#%{HOST}#$HOST#g" ./swaggerui/swagger.json
    # Replace default URL with latest commit ID to avoid browser caching
    sed -i '' -e "s|url:.*|url: \"./swagger.json?$(git rev-parse --short HEAD)\",|" ./swaggerui/index.html
else
    # Replace HOST by corresponding env var
    sed -i -e "s#%{HOST}#$HOST#g" ./swaggerui/swagger.json
    # Replace default URL with latest commit ID to avoid browser caching
    sed -i -e "s|url:.*|url: \"./swagger.json?$(git rev-parse --short HEAD)\",|" ./swaggerui/index.html
fi


