#!/usr/bin/env bash

# 初始化数据库
user="root"
password="123456"
db="cloudops"
prot="3306"
host="127.0.0.1"

# 创建数据库
mysql -u$user -p$password -h$host -P$prot <<EOF
CREATE DATABASE IF NOT EXISTS $db DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
EOF
# 导入数据
mysql -u$user -p$password -h$host -P$prot $db < ./sql/cloudops-init.sql
