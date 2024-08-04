#!/bin/sh

sed -i "s~<ENV>~${ENV}~g" /usr/share/nginx/html/config.js
sed -i "s~<API_BASE_URL>~${API_BASE_URL}~g" /usr/share/nginx/html/config.js
sed -i "s~<AUTH_BASE_URL>~${AUTH_BASE_URL}~g" /usr/share/nginx/html/config.js
sed -i "s~<AUTH_CLIENT_ID>~${AUTH_CLIENT_ID}~g" /usr/share/nginx/html/config.js
sed -i "s~<LEADERBOARD_TYPE>~${LEADERBOARD_TYPE}~g" /usr/share/nginx/html/config.js
sed -i "s~<SLACK_CHANNEL_ID>~${SLACK_CHANNEL_ID}~g" /usr/share/nginx/html/config.js
exec "$@"
