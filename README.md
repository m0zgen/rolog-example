# RoLog Example

This is a demo app for usage Rotate Log module.

## Parameters

- `logDir` - Log store catalog
- `staticFilename` - Permanent log file name
- `archivePattern` - Backup log file name
- `zippedArchive` - Pack backup log to `zip` format
- `maxSize` - Max log file size in megabytes
- `maxAge` - Max log age in days
- `maxBackups` - Number of backup logs
- `checkInterval` - Check backup need in hours
- `bufferSize` - Size of the log message buffer
- `logLevel` - Set the log level to Info

## Usage Example

```go
logDir := "logs"
staticFilename := "application.log"     // Permanent log file name
archivePattern := "application-%s.log"  // Backup log file name
zippedArchive := true                   // Archive log yes/no
maxSize := 1                            // 20 MB
maxAge := 14                            // 14 days
maxBackups := 3                         // 3 backups
checkInterval := time.Hour              // Check every hour
bufferSize := 100                       // Size of the log message buffer
logLevel := logrus.InfoLevel            // Set the log level to Info
```