# `logfrog`


**`logfrog` likes json logs and it helps you to like them too**

`logfrog` is a cli program that processes json logs line by line from stdin. Those logs typically come from loggers like [logrus](https://github.com/Sirupsen/logrus), [zap](https://github.com/uber-go/zap), [apex](https://github.com/apex/log) or [others](https://github.com/topics/structured-logging), that let you write logs as json objects.

Status: logfrog is very young atm and especially the way we filter is most likely going to change. Despite that it already provides a lot of value, when you are trying to make sense of logs.

## installation

`brew install foomo/logfrog/logfrog`

## use cases

### stern

### docker-compose logs

```bash
stern -o json -n some-name-space | logfrog -log-type stern
```

```bash
docker-compose logs --no-color -f | logfrog -log-type docker-compose
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

// filter function must be named filter, it will be reloaded if updated
//
// @param logEntry:{msg?:string;level?:string;time?:string, ...}
// @param service:string only set with -log-type docker-compose or stern
//
// @return logEntry | null when null is returned this entry is filtered out
function filter(logEntry, service) {
  // let us look at the service in this naive docker-compose example I butcher the name
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

- ☑️ map more fields
- ☑️ maybe add a web frontend ?!
- ✅ stern mode like docker-compose
- ✅ add hombrew support
