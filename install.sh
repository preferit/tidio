#!/bin/bash -e
systemctl stop tidio
cp tidio ../bin/tidio
systemctl start tidio

cp etc/tidio.preferit.se /etc/nginx/sites-available/
systemctl reload nginx
