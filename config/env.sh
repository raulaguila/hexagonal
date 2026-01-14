#!/bin/bash

ipaddr=$(hostname -I | cut -d' ' -f1)
release_version=$(cat configs/version.txt | tr -d "[:space:]")
access_token=$(openssl genrsa 2048 | base64 | tr -d \\n)
refresh_token=$(openssl genrsa 2048 | base64 | tr -d \\n)

echo "TZ='America/Manaus'                             # Set system time zone
SYS_VERSION='${release_version}'                    # System version
ENVIRONMENT='production'                       # System environment

LOG_FORMAT='json'                               # Log format
LOG_LEVEL='info'                                # Log level

API_PORT='9000'                                 # API Container PORT
API_LOGGER='1'                                  # API Logger enable
API_SWAGGO='1'                                  # API Swagger enable
API_CACHE='0'                                   # API Cache enable
API_ENABLE_PREFORK='1'                          # API enable fiber prefork
API_DEFAULT_SORT='updated_at'                   # API default column sort
API_DEFAULT_ORDER='desc'                        # API default order
API_ACCEPT_SKIP_AUTH='1'                        # API accept skip auth header

ACCESS_TOKEN_EXPIRE='15'                        # Access token expiration time in minutes
RFRESH_TOKEN_EXPIRE='60'                        # Refresh token expiration time in minutes

ACCESS_TOKEN='${access_token}'                  # Token to encode access token - PRIVATE TOKEN
RFRESH_TOKEN='${refresh_token}'                 # Token to encode refresh token - PRIVATE TOKEN

POSTGRES_HOST='postgres'                        # Postgres Container HOST
POSTGRES_PORT='5432'                            # Postgres Container PORT
POSTGRES_USER='root'                            # Postgres USER
POSTGRES_PASS='root'                            # Postgres PASS
POSTGRES_BASE='api'                             # Postgres BASE

MINIO_HOST='\${ipaddr}'                          # Minio HOST
MINIO_API_PORT='9004'                           # Minio API PORT
MINIO_WEB_PORT='9005'                           # Minio WEB PORT
MINIO_USER='minio'                              # Minio USER
MINIO_PASS='miniopass'                          # Minio PASS
MINIO_BUCKET_FILES='api'                        # Minio BUCKET" >.env
