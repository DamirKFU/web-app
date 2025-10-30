mkdir -p /run/postgresql
chown postgres:postgres /run/postgresql
su - postgres -c "pg_ctl -D /var/lib/postgresql/data -l /var/lib/postgresql/logfile start"
redis-server &