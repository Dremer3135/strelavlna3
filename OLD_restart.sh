git pull
systemctl stop valkey.service
rm dump.rdb
cp dump2.rdb dump.rdb
systemctl start valkey.service
systemctl restart enigine.service
journalctl -xeu enigine.service -f
