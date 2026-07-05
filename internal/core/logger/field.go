package core_logger

import "time"

// Field — структурированное поле лога.
type Field struct {
	key string
	val any
}

func String(key, val string) Field {
	return Field{key: key, val: val}
}

func Int(key string, val int) Field {
	return Field{key: key, val: val}
}

func Error(err error) Field {
	return Field{key: "error", val: err}
}

func Time(key string, val time.Time) Field {
	return Field{key: key, val: val}
}

func Duration(key string, val time.Duration) Field {
	return Field{key: key, val: val}
}
