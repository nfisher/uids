#!/bin/sh -eu

yum update -y
yum install -y mariadb-server
systemctl start mariadb.service
systemctl enable mariadb.service
systemctl stop firewalld
systemctl disable firewalld.service
mysql -u root -h localhost <<EOT
GRANT ALL ON *.* TO root@'%' IDENTIFIED BY 'secret123';
EOT

/etc/init.d/vboxadd setup
