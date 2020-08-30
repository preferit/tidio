#!/bin/bash -e

# tidio service
systemctl stop tidio
cp tidio ../bin/tidio
systemctl start tidio

# nginx configuration
cp etc/tidio.preferit.se /etc/nginx/sites-available/
systemctl reload nginx
