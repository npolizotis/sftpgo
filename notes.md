# Debug
```bash
docker run -it --rm --entrypoint /bin/sh \
-e R2_BUCKET="${R2_BUCKET}" \
-e R2_ACCESS_KEY="${R2_ACCESS_KEY}" \
-e R2_SECRET_KEY="${R2_SECRET_KEY}" \
-e R2_DATA_PATH="${R2_DATA_PATH}" \
-e R2_URL="${R2_URL}" \
npolizotis:sftpgo ```

# Export R2 variables

```bash
set -a 
. dev-r2.rel
set +a
```