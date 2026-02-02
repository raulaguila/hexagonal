#!/bin/bash

ipaddr=$(hostname -I | cut -d' ' -f1)
release_version=$(cat config/version.txt | tr -d "[:space:]")
access_token=$(openssl genrsa 2048 | base64 | tr -d \\n)
refresh_token=$(openssl genrsa 2048 | base64 | tr -d \\n)

echo "TZ='America/Manaus'                             # Set system time zone
SYS_NAME='backend'                              # System name
SYS_VERSION='${release_version}'                             # System version
SYS_ENVIRONMENT='production'                    # System environment

LOG_FORMAT='json'                               # Log format
LOG_LEVEL='info'                                # Log level (debug, info, warn, error, fatal, default=info)

API_PORT='9999'                                 # API Container PORT
API_LOGGER='1'                                  # API Logger enable
API_SWAGGO='1'                                  # API Swagger enable
API_CACHE='0'                                   # API Cache enable
API_PREFORK='0'                                 # API enable fiber prefork
API_DEFAULT_SORT='updated_at'                   # API default column sort
API_DEFAULT_ORDER='desc'                        # API default order
API_ACCEPT_SKIP_AUTH='1'                        # API accept skip auth header

ACCESS_TOKEN_EXPIRE='50m'                       # Access token expiration (m=min, s=seg, h=hour, default=50m)
RFRESH_TOKEN_EXPIRE='3h'                        # Refresh token expiration (m=min, s=seg, h=hour, default=3h)

ACCESS_TOKEN='${access_token}'                  # Token to encode access token - PRIVATE TOKEN
RFRESH_TOKEN='${refresh_token}'                 # Token to encode refresh token - PRIVATE TOKEN

POSTGRES_HOST='postgres'                        # Postgres Container HOST
POSTGRES_PORT='5432'                            # Postgres Container PORT
POSTGRES_USER='root'                            # Postgres USER
POSTGRES_PASS='root'                            # Postgres PASS
POSTGRES_BASE='api'                             # Postgres BASE

REDIS_HOST='redis'                              # Redis HOST
REDIS_PORT='6379'                               # Redis PORT
REDIS_USER='default'                            # Redis USER
REDIS_PASS='redispass'                          # Redis PASS
REDIS_DB='0'                                    # Redis DB
REDIS_TTL='10m'                                 # Redis TTL

MINIO_HOST='${ipaddr}'                          # Minio HOST
MINIO_API_PORT='9004'                           # Minio API PORT
MINIO_WEB_PORT='9005'                           # Minio WEB PORT
MINIO_USER='minio'                              # Minio USER
MINIO_PASS='miniopass'                          # Minio PASS
MINIO_BUCKET_FILES='api'                        # Minio BUCKET

ELASTICSEARCH_PORT='9200'                       # Elasticsearch PORT
ELASTICSEARCH_HOST='127.0.0.1'                  # Elasticsearch HOST
ELASTIC_USER='elastic'                          # Elastic USER
ELASTIC_PASS='changeme'                         # Elastic PASS

KIBANA_PORT='5601'                              # Kibana PORT
KIBANA_SYSTEM_USER='kibana_system'              # Kibana System USER
KIBANA_SYSTEM_PASS='changeme'                   # Kibana System PASS

APM_PORT='8200'                                 # APM Server PORT
OTLP_GRPC_PORT='4317'                           # OTLP gRPC PORT
OTLP_HTTP_PORT='4318'                           # OTLP HTTP PORT" >.env
