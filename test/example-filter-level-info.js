function filter(service, logEntry) {
    if(logEntry.level === "info") {
        return null;
    }
    return logEntry;
}