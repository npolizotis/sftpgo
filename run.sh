#!/bin/sh
set -e

echo "Check diskspace on VM"
df -h

echo "litestream version"
litestream version

echo "Restore db if exists"
set +e
litestream_restore=$(litestream restore -if-replica-exists /var/lib/sftpgo/sftpgo.db )
if  [ $? -ne 0 ]; then
	# Delete and retry
	rm /var/lib/sftpgo/sftpgo.db
	litestream restore -if-replica-exists /var/lib/sftpgo/sftpgo.db
	echo "Restored successfully"
else
	echo "Restore not applied"
fi
#litestream restore -if-replica-exists /pb_data/logs.db
set -e

echo "replicate!"
exec litestream replicate -exec "/usr/local/bin/sftpgo serve"
