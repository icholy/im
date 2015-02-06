# IM (WIP)

> An tool for keeping track of what I'm doing

## Todo

- [x] Basic Day & Task storage with tags
- [x] Update date with ping
- [x] Global locking on data directory
- [ ] Query interface
- [ ] Web based calendar view

## Ping

``` sh
$ crontab -e
```

Add the following

```
*/5 * * * * /usr/local/bin/im -ping
```

## Submit a task

``` sh
$ im doing some @stuff
```




