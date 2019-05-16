#!/bin/bash

cat <<-EOF >>~/.bashrc
	echo " ____                      _                 _           "
	echo "|  _ \\\\  _____      ___ __ | | ___   __ _  __| | ___ _ __ "
	echo "| | | |/ _ \\\\ \\\\ /\\\\ / / '_ \\\\| |/ _ \\\\ / _\\\` |/ _\\\` |/ _ \\\\ '__|"
	echo "| |_| | (_) \\\\ V  V /| | | | | (_) | (_| | (_| |  __/ |   "
	echo "|____/ \\\\___/ \\\\_/\\\\_/ |_| |_|_|\\\\___/ \\\\__,_|\\\\__,_|\\\\___|_|   "
	echo "                                                         "
	echo ""
	echo -e "when running with DOCKER_TEST=true, try:\n"
	echo -e "    \\\$ curl -d \"uri=config-1.tar.gz&unarchive=true\" localhost:9000/v1/download\n"
	echo -e "    using dummy files in /tmp (bucket) and /tmp/downloads (download dir)"
	echo ""
EOF

if [ "$DOCKER_TEST" = true ]; then
	LOG_LEVEL="debug"
	BUCKET_PROTO="local"
	BUCKET_NAME="/tmp"
	DOWNLOAD_DIR="/tmp/downloads"
	KEEP_OLD_COUNT=2

	# Creates dummy config files for downloader
	cd ${BUCKET_NAME}
	for a in {1..10}; do
		temp_config="config-$a"

		touch "${temp_config}.yaml"
		tar cvzf "${temp_config}.tar.gz" "${temp_config}.yaml"
		rm -rf "${temp_config}.yaml"
	done
fi

echo "starting configdownloader..."
exec /bin/configdownloader \
	-logLevel ${LOG_LEVEL} \
	-bucketProto "${BUCKET_PROTO}" \
	-bucketName "${BUCKET_NAME}" \
	-downloadDIR "${DOWNLOAD_DIR}" \
	-keepOldCount ${KEEP_OLD_COUNT}
