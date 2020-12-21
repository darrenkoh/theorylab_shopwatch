#!/bin/bash

# Docker Reference - https://neo4j.com/docs/operations-manual/3.5/docker/#docker
# Neo4J GraphGist - https://neo4j.com/graphgists/
# Neo4J Desktop App - https://neo4j.com/download/?ref=developer-neo4j-desktop


mkdir -p neo4jdata
mkdir -p neo4jlog

docker run \
    --name neo4j \
    --publish=7474:7474 \
    --publish=7687:7687 \
    --volume=$PWD/neo4jdata:/data \
    --volume=$PWD/neo4jlog:/logs \
    --env NEO4J_dbms_memory_pagecache_size=1G \
    --env NEO4J_AUTH=neo4j/masterxzy \
    --env=NEO4J_ACCEPT_LICENSE_AGREEMENT=yes \
    neo4j:enterprise