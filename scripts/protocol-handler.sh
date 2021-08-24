#!/usr/bin/env sh

set -euf

decrypt_command="$HELM_PLUGIN_DIR/bin/helm-gsm-linux -f"

# https://helm.sh/docs/topics/plugins/#downloader-plugins
# It's always the 4th parameter
file=$(printf '%s' "${4}" | sed -E -e 's!gsm://!!')

# send output to /dev/null so it doesn't break helm
$decrypt_command ${file} > /dev/null
real_file="${file}.dec"

cat "${real_file}"
