# weibo监听

Run example

```shell
go run main.go -u 1223178222 --dsn="root:root@/weibo"
```

| Flags         | description                     |
|:--------------|:--------------------------------|
| -c / --cookie | weibo cookie                    |
| -u / --userid | weibo uesr id                   |
| -p / --page   | start page                      |
| -s / --sleep  | request interval                |
| -d / --dn     | database driver, mysql          |
| --dsn         | database connection information |
| -f / --full   | crawl all weibo                 |
| --cron        | cron rlus                       |
| -t / --tz     | time zone                       |