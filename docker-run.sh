#!/bin/bash
set -xauo

echo "Bucket is: $R2_BUCKET"
docker build --build-arg FEATURES="nos3,nogcs" -t npolizotis/sftpgo -f Dockerfile.alpine . && \
docker run -it --rm --publish 8080:8080 --publish 2022:10022 \
-e R2_BUCKET="${R2_BUCKET}" \
-e R2_ACCESS_KEY="${R2_ACCESS_KEY}" \
-e R2_SECRET_KEY="${R2_SECRET_KEY}" \
-e R2_DATA_PATH="${R2_DATA_PATH}" \
-e R2_URL="${R2_URL}" \
npolizotis/sftpgo 
