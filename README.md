# Doing (WIP)

> An tool for keeping track of what I've been doing

## What

* Keeps track of when your computer was turned on and off.
* Can run a webserver with a calendar view.
* Makes you look busier than you actually are.
* Treats words prefixed with `@` as tags.
* Tab completes existing tags.

## UI

Entering Records:

``` sh
$ doing <description>
$ doing @app refactoring Foo controller
$ doing listening to @dubstep on @youtube
$ doing nothing
```

Starting Web Interface:

``` sh
$ doing --web
```

Updating Workday:

``` sh
$ doing --ping
```

## How

* A single entry is called a `Record`.
* A startup and shutdown cycle is called a `Workday`.
* A cron job periodically calls `doing --ping` to update the `Workday`.
* Storage will be implemented using git.


