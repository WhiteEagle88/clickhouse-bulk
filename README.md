# ClickHouse-Bulk

Simple [Yandex ClickHouse](https://clickhouse.yandex/) insert collector. It collect requests and send to ClickHouse servers.

### Features
- Group n requests and send to any of ClickHouse server
- Sending collected data by interval
- Tested with VALUES, TabSeparated formats
- Supports many servers to send
- Supports query in query parameters and in body
- Supports other query parameters like username, password, database
- Supports basic authentication
- Deduplication of queries using bigcache
 

For example:
```sql
INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3')
INSERT INTO table3 (c1, c2, c3) VALUES ('v4', 'v5', 'v6')
```
sends as
```sql
INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3')('v4', 'v5', 'v6')
```


### Options
- -config - config file (json); default _config.json_


### Configuration file
```json5
{
  "listen": ":8124", 
  "flush_count": 10000, // check by \n char
  "flush_interval": 1000, // milliseconds
  "dump_check_interval": 300, // interval for try to send dumps (seconds); -1 to disable
  "debug": false, // log incoming requests
  "dump_dir": "dumps", // directory for dump unsended data (if clickhouse errors)
  "clickhouse": {
    "down_timeout": 60, // wait if server in down (seconds)
    "connect_timeout": 10, // wait for server connect (seconds)
    "servers": [
      "http://127.0.0.1:8123"
    ],
  },
  "cache": {
    "shards": 1024, // Number of cache shards, value must be a power of two
    "life_window": 10, // Time after which entry can be evicted(in minutes)
    "clean_window": 0, // Interval between removing expired entries (clean up). If set to <= 0 then no action is performed.
    "max_entries_in_window": 600000, // Max number of entries in life window. Used only to calculate initial size for cache shards.
    "max_entry_size": 500, // Max size of entry in bytes. Used only to calculate initial size for cache shards.
    "verbose": true, // Verbose mode prints information about new memory allocation
    "max_cache_size": 0 // Limit for cache size in MB. Cache will not allocate more memory than this limit.
  }
}
```

### Environment variables (used for docker image)

* `CLICKHOUSE_SERVERS` - comma separated list of servers
* `CLICKHOUSE_FLUSH_COUNT` - count of rows for insert
* `CLICKHOUSE_FLUSH_INTERVAL` - insert interval
* `DUMP_CHECK_INTERVAL` - interval of resend dumps  
* `CLICKHOUSE_DOWN_TIMEOUT` - wait time if server is down  
* `CLICKHOUSE_CONNECT_TIMEOUT` - clickhouse server connect timeout  

### Quickstart

`./clickhouse-bulk`
and send queries to :8124

### Metrics
manual check main metrics
`curl -s http://127.0.0.1:8124/metrics | grep "^ch_"`
* `ch_bad_servers 0` - actual count of bad servers
* `ch_dump_count 0` - dumps saved from launch
* `ch_queued_dumps 0` - actual dump files id directory
* `ch_good_servers 1` - actual good servers count
* `ch_received_count 40` - received requests count from launch
* `ch_sent_count 1` - sent request count from launch


### Tips

For better performance words FORMAT and VALUES must be uppercase.

### Credits

This repository is based on the following work:
- Nikolay's GitHub Project [nikepan/clickhouse-bulk](https://github.com/nikepan/clickhouse-bulk). Thanks a lot for his work.