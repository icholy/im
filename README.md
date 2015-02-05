# Doing (WIP)

> An tool for keeping track of what I've been doing

## What

* Keeps track of when your computer was turned on and off.
* Can run a webserver with a calendar view.
* Makes you look busier than you actually are.
* Treats words prefixed with `@` as tags.
* Tab completes existing tags.

## UI

``` sh
$ doing <description>
$ doing @app refactoring Foo controller
$ doing listening to @dubstep on @youtube
$ doing nothing
```

## How

* A single entry is called a `Record`.
* A daemon runs in the backround and monitors system uptime.
* The daemon exposes a rest interface to manage records.
* A startup and shutdown cycle is called a `Workday`.
* Not sure how storage is going to be implemented.


