/*

ClickHouse-Bulk

Simple Yandex ClickHouse (https://clickhouse.yandex/) insert collector. It collect requests and send to ClickHouse servers.


Features

- Group n requests and send to any of ClickHouse server

- Sending collected data by interval

- Tested with VALUES, TabSeparated formats

- Supports many servers to send

- Supports query in query parameters and in body

- Supports other query parameters like username, password, database

- - Supports basic authentication


For example:

INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3')

INSERT INTO table3 (c1, c2, c3) VALUES ('v4', 'v5', 'v6')

sends as

INSERT INTO table3 (c1, c2, c3) VALUES ('v1', 'v2', 'v3')('v4', 'v5', 'v6')


Options

- -config - config file (json); default _config.json_


Configuration file

{
  "listen": ":8124",
  "flush_count": 10000, // check by \n char
  "flush_interval": 1000, // milliseconds
  "debug": false, // log incoming requests
  "dump_dir": "dumps", // directory for dump unsended data (if clickhouse errors)
  "clickhouse": {
    "down_timeout": 300, // wait if server in down (seconds)
    "servers": [
      "http://127.0.0.1:8123"
    ]
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


Quickstart

`./clickhouse-bulk`
and send queries to :8124


Tips

For better performance words FORMAT and VALUES must be uppercase.

*/
package main
