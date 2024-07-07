function filter(logEntry, service) {
    if(logEntry.level === "info") {
        return null;
    }
    return logEntry;
}
