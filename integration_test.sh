#!/bin/bash
set +e

# All environment variable
PGDATABASE=${PGDATABASE}
PGUSER=${PGUSER}
PGHOST=${PGHOST}
PGPORT=${PGPORT}
PGPASSWORD=${PGPASSWORD}
MOCK_DATA_TEST_RUNNER=${MOCK_DATA_TEST_RUNNER}

# Logger
function logger() {
    # Priority & Message Extraction
    local log_priority=$1
    local log_message=$2

    # All the color module
    local RED='\033[0;31m'       # Error
    local NC='\033[0m'           # No Color
    local YELLOW='\033[1;33m'    # Warning
    local BLUE='\033[0;34m'      # Info
    local PURPLE='\033[0;35m'    # Debug
    local GREEN='\033[0;32m'     # UnKnown Priority

    # Get the verbosity
    local verbose_mode=${VERBOSE}

    # Print the logger message
    case ${log_priority} in
        "INFO")
          echo -e "${BLUE}[$(date +"%Y-%m-%d %H:%M:%S")]:${log_priority}${NC}:${log_message}"
          ;;
        "ERROR")
          echo -e "${RED}[$(date +"%Y-%m-%d %H:%M:%S")]:${log_priority}${NC}:${log_message}"
          ;;
        "DEBUG")
          if [ "${verbose_mode}" == "YES" ]; then
             echo -e "${PURPLE}[$(date +"%Y-%m-%d %H:%M:%S")]:${log_priority}${NC}:${log_message}"
          fi
          ;;
        "WARNING")
          echo -e "${YELLOW}[$(date +"%Y-%m-%d %H:%M:%S")]:${log_priority}${NC}:${log_message}"
          ;;
        *)
          echo -e "${GREEN}[$(date +"%Y-%m-%d %H:%M:%S")]:${log_priority}${NC}:${log_message}"
          ;;
    esac
}

# Check if the command was a success or failure
function success_or_failure() {
    retcode=$1; message=$2;
    if [ "${retcode}" == 0 ]; then
        logger "INFO" "${message} was a success..."
    else
        logger "ERROR" "${message} was a failure..."
        exit 255
    fi
}
logger "INFO" "Starting the integration tests"

### Database command
# Creating database objects
logger "INFO" "TEST: Creating demo database"
go run . d -c
success_or_failure $? "Creating demo database"

# Mocking the whole database
logger "INFO" "TEST: Mocking demo database"
go run . d -f -q
success_or_failure $? "Mocking the whole demo database"

# Create and mock the database at the same time
logger "INFO" "TEST: Creating and mocking database"
go run . d -f -q -c
success_or_failure $? "Creating and mocking database"

### Tables command
# Creating fake tables
logger "INFO" "TEST: Creating fake tables"
go run . t -c -n 30 -m 20 -j -x "mock-data-test-tables" -y "mock-data-test-col" -s "public" -q
success_or_failure $? "Creating fake tables"

# Mock tables from the list
logger "INFO" "TEST: Mocking selected tables"
go run . t -t actor,film -s "public" -q
success_or_failure $? "Mocking selected tables"

### Schema command
# Mocking all tables of the schema
logger "INFO" "TEST: Mocking tables of schema"
go run . s -n "public" -q
success_or_failure $? "Mocking tables of schema"

### Custom command
# Create a custom plan for tables
logger "INFO" "TEST: Creating custom loading plan"
go run . c -t actor -q
success_or_failure $? "Creating custom loading plan"

# Loading the data based on custom files
logger "INFO" "TEST: Loading from custom generated file"
cat > custom.yml << EOF
Custom:
- Schema: public
  Table: actor
  Column:
  - Name: first_name
    Type: character varying(45)
    UserData: []
    Realistic: "NameFirstName"
  - Name: last_name
    Type: character varying(45)
    UserData: ["Murphy", "Soil", "Hu", "Yu"]
    Realistic: ""
  - Name: last_update
    Type: timestamp without time zone
    UserData: []
    Realistic: ""
EOF
go run . c -f custom.yml -q
success_or_failure $? "Mocking tables using custom file"


logger "INFO" "Completed successfully all the integration tests"