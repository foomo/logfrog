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
