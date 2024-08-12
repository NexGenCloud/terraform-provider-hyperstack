#!/bin/bash

# Global variables for configuration
URL="http://${SERVER_IP:-localhost}/v1/completions"
CONTENT_TYPE="Content-Type: application/json"
MODEL="${MODEL_NAME:-NousResearch/Meta-Llama-3-8B-Instruct}"
PROMPT="Kubernetes with GPU is a"
MAX_TOKENS=1000
TEMPERATURE=0
CLUSTER="${CLUSTER_NAME:-undefined}"

# Function to perform a CURL request to the endpoint
perform_curl_request() {
    local response
    response=$(curl -s -w "\n%{http_code}" -H "$CONTENT_TYPE" -d '{
        "model": "'"$MODEL"'",
        "prompt": "'"$PROMPT"'",
        "max_tokens": '"$MAX_TOKENS"',
        "temperature": '"$TEMPERATURE"'
    }' "$URL")

    if [[ -z $response ]]; then
        printf "Error: No response from server.\n" >&2
        return 1
    fi

    # Split response and HTTP status code
    local http_status
    http_status=$(tail -n1 <<< "$response")
    response=$(sed '$ d' <<< "$response")

    # Check HTTP status code
    if [[ "$http_status" -ne 200 ]]; then
        printf "Error: HTTP request failed with status code %s\n" "$http_status" >&2
        return 1
    fi

    # Validate the JSON response
    if ! validate_json_response "$response"; then
        printf "Error: Response is not a valid JSON\n" >&2
        return 1
    fi

    # Parse the JSON response
    parse_json_response "$response"
}

# Function to validate JSON response
validate_json_response() {
    local json="$1"
    if ! jq -e . >/dev/null 2>&1 <<<"$json"; then
        return 1
    fi
    return 0
}

# Function to parse the JSON response
parse_json_response() {
    local json="$1"
    local id
    local model
    local text
    id=$(jq -r '.id' <<<"$json")
    model=$(jq -r '.model' <<<"$json")
    text=$(jq -r '.choices[0].text' <<<"$json")

    if [[ -n $id && -n $model && -n $text ]]; then
        printf "Cluster name: %s\n" "$CLUSTER"
        printf "ID: %s\nModel: %s\nText: %s\n" "$id" "$model" "$text"
    else
        printf "Error: Missing fields in the response.\n" >&2
        return 1
    fi
}

# Main function
main() {
    if ! perform_curl_request; then
        printf "Failed to complete the CURL request.\n" >&2
        exit 1
    fi
}

main "$@"
