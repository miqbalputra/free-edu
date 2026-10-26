#!/bin/bash
# Backup harian SIMAKA - rotasi 7 hari
DIR=/home/agentuser/backups/simaka
sudo -u postgres pg_dump simaka | gzip > "$DIR/simaka-$(date +%F).sql.gz"
find "$DIR" -name 'simaka-*.sql.gz' -mtime +7 -delete
