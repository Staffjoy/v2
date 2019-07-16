#!/bin/bash
set -e

sudo apt install -y -q  mysql-client

echo -e -n "Waiting for mysql to finish booting..."
until timeout 1 nc -z 10.0.0.100 3306 &>/dev/null;
do
    echo -n "."
    sleep 1
done

echo "MySQL UP - Initializing databases"

# account-mysql-service
mysql -u root -pSHIBBOLETH -h 10.0.0.100 -e "create database account"
mysql -u root -pSHIBBOLETH -h 10.0.0.100 -e "create database company"