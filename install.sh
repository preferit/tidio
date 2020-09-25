#!/bin/bash -e

# tidio service
systemctl stop tidio
cp systemd.service /lib/systemd/system/tidio.service
systemctl daemon-reload
cp tidio ../bin/tidio
systemctl start tidio

# nginx configuration
cp nginx.conf /etc/nginx/sites-available/tidio.preferit.se
systemctl reload nginx
