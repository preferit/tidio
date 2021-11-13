#!/bin/bash -e

# This script assumes the working directory is /tmp/tidio

# tidio service
systemctl stop tidio
cp systemd.service /lib/systemd/system/tidio.service
systemctl daemon-reload
cp tidio /usr/local/bin/tidio
systemctl start tidio

# nginx configuration
cp nginx.conf /etc/nginx/sites-available/tidio.preferit.se
systemctl reload nginx
