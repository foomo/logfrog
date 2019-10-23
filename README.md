# `logfrog`

**`logfrog` likes json logs and it helps you to like them too**

`logfrog` is a cli program that processes json logs line by line from stdin. Those logs typically come from loggers like [logrus](https://github.com/Sirupsen/logrus), [zap](https://github.com/uber-go/zap), [apex](https://github.com/apex/log) or [others](https://github.com/topics/structured-logging), that let you write logs as json objects.

## use cases

### docker-compose logs

```bash
docker-compose logs --no-color -f | logfrog --docker-compose
```

### docker

```bash
docker logs some-container 2>&1 | logfrog
```

### json log files

```bash
tail -f path-to-file.json | logfrog
```

## js filtering

`logfrog` lets you transform and filter json logs with a javascript function named `filter` that must be defined in a .js file that is passed with the ar `--js-filter`

```bash
tail -f path-to-file.json | logfrog --js-filter path/to/filter.js
```

- the js file is executed with the otto vm [https://github.com/robertkrimen/otto](https://github.com/robertkrimen/otto)

- it has to contain a filter function like the one below

- the file will be reevaluated, when it changes

- *this is highly EXPERIMENTAL* and we would love to hear back from you

  

```JavaScript
// filter function must be named filter
//
// @param service:string only filled with logfrog --docker-compose
// @param logEntry:{msg?:string;level?:string;time?:string, ...}
// @return logEntry | null when null is returned this entry is filtered out
function filter(service, logEntry) {
  // let us look at the service
  switch (service.substr(0, service.length - 2)) {
    case "elasticsearch":
      // very minimal log entries for elastic search
      return { level: logEntry.level, msg: logEntry.message };
  }
  // log entry manipulation
  // some date formatting
  logEntry.msg = logEntry.msg.substr(0, 256);
  logEntry.time = new Date(logEntry.time).toLocaleString();
  // trimming a stack
  if (logEntry.stack) {
    logEntry.stack = logEntry.stack.substr(0, 300) + " ...";
  }
  // go crazy and have fun ;)
  return logEntry;
}

```

## standard fields

This is an initial set of fields, please let us know what we should add.

- `msg` <- msg, message, Message
- `level` <- level, Level
- `time` <- time, timestamp, Timestamp
  

## todos

- add hombrew support
- map more fields
- stern mode like docker-compose
- maybe add a web frontend ?!