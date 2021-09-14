#!/usr/bin/env sh

set -euf

os=$(uname)
binary="helm-gsm-linux"
if [ "$os" = "Darwin" ]; then
  binary="helm-gsm-darwin"
fi
decrypt_command="$HELM_PLUGIN_DIR/bin/$binary -b -f"

# https://helm.sh/docs/topics/plugins/#downloader-plugins
# It's always the 4th parameter
file=$(printf '%s' "${4}" | sed -E -e 's!gsm://!!')

# send output to /dev/null so it doesn't break helm
$decrypt_command ${file} > /dev/null
real_file="${file}.dec"

cat "${real_file}"
