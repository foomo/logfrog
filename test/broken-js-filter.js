functio filter(logEntry) {
    if(logEntry.level === "info") {
        return null;
    }
    return logEntry;
}