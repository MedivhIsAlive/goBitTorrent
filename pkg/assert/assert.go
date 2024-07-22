package assert

import (
	"log"
	"log/slog"
)

var assertData map[string]any = map[string]any{}

func AddAssertData(key string, value any) {
	assertData[key] = value
}

func RemoveAssertData(key string) {
	delete(assertData, key)
}

func logAssertFailed(msg string) {
	for k, v := range assertData {
		slog.Error("context", "key", k, "value", v)
	}
	log.Fatal(msg)
}

func Assert(truth bool, msg string) {
	if !truth {
		logAssertFailed(msg)
	}
}

func NoError(err error, msg string) {
	if err != nil {
		slog.Error("NO error #error encountered", "error", err)
		logAssertFailed(msg)
	}
}
