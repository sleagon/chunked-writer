# Chunked Writer


## Notice

[```lumberjack```][1] is nice rotate logger, just use it.


## Purpose

The only purpose of this package is reduce the complexity of splitting logs. This package is extremely simple.


## Usage

This writer can be used with zap/zerolog or any logger can config writer.

```go
// work with zerolog
func main() {
	w, err := writer.New("", "", 20)
	if err != nil {
		panic(err)
	}
	Logger := zerolog.New(w).With().Timestamp().Logger()
	Logger.Print("good")
}

// new writer manually
	w := &Writer{
		ChunkSize: DefaultChunkSize,
		Prefix:    DefaultPrefix,
		Pattern:   DefaultPattern,
		Dir:       DefaultDir,
    }
```



[1]: https://github.com/natefinch/lumberjack