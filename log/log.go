package log

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/natefinch/lumberjack"
)

var (
	defaultFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	defaultLevel  = "debug"
	defaultOutput = os.Stderr

	log = logrus.New()
)

func init() {
	SetFormatter(defaultFormatter)
	SetLevel(defaultLevel)
	SetOutput(defaultOutput)
	AddHook(&CallerHook{})

	c := make(chan os.Signal, 1)

	//syscall.SIGUSR1,syscall.SIGUSR2
	signal.Notify(c, syscall.Signal(0x1e), syscall.Signal(0x1f))
	go watchAndUpdateLoglevel(c, log)
}

func watchAndUpdateLoglevel(c chan os.Signal, logger *logrus.Logger) {
	for {
		select {
		case sig := <-c:
			if sig == syscall.Signal(0x1e) {
				level := logger.Level
				if level == logrus.PanicLevel {
					fmt.Println("Raise log level: It has been already the most top log level: panic level")
				} else {
					logger.Level = level - 1
					fmt.Println("Raise log level: the current level is", logger.Level)
				}

			} else if sig == syscall.Signal(0x1f) {
				level := logger.Level
				if level == logrus.DebugLevel {
					fmt.Println("Reduce log level: It has been already the lowest log level: debug level")
				} else {
					logger.Level = level + 1
					fmt.Println("Reduce log level: the current level is", logger.Level)
				}

			} else {
				fmt.Println("receive unknown signal:", sig)
			}
		}
	}
}

// SetFormatter set log formatter.
func SetFormatter(formatter logrus.Formatter) {
	log.Formatter = formatter
}

// SetLevel set log level.
func SetLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return
	}

	log.Level = level
}

// AddHook adds a hook to the logger hooks.
func AddHook(hook logrus.Hook) {
	log.Hooks.Add(hook)
}

// SetOutput set log output.
func SetOutput(out io.Writer) {
	log.Out = out
}

// SetOutputLogFile set output log file
func SetOutputLogFile(filepath string) {
	log.Out = &lumberjack.Logger{
		Filename:   filepath, //"./app.log",
		MaxSize:    1,        // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}
}

// Panic log a panicf message
func Panic(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Panicf(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Panic(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Panicf(format, v...)
	}
}

// Fatal log a fatalf message
func Fatal(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Fatalf(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Fatal(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Fatalf(format, v...)
	}
}

// Error log a error message
func Error(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Errorf(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Error(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Errorf(format, v...)
	}
}

// Warn log a warn message
func Warn(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Warnf(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Warn(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Warnf(format, v...)
	}
}

// Info log a info message
func Info(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Infof(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Info(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Infof(format, v...)
	}
}

// Debug log a debug message
func Debug(f interface{}, v ...interface{}) {
	var format string
	switch f.(type) {
	case string:
		format = f.(string)
		log.Debugf(format, v...)
	default:
		format = fmt.Sprint(f)
		if len(v) == 0 {
			log.Debug(format)
			return
		}
		format += strings.Repeat(" %v", len(v))
		log.Debugf(format, v...)
	}
}
