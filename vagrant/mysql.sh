#!/bin/bash
set -e

sudo apt install -y -q  mysql-client

# Run mysql base provisioning
sleep 120 # to let it boot

echo "Initializing databases"

# account-mysql-service
mysql -u root -pSHIBBOLETH -h 10.0.0.100 -e "create database account"
mysql -u root -pSHIBBOLETH -h 10.0.0.100 -e "create database company"