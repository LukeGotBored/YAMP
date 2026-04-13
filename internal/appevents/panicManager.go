package appevents

import (
        "fmt"
        "time"
        "github.com/wailsapp/wails/v3/pkg/application"
)

const (
        BackendFatalErrorEvent = "app:fatal-error"
        BackendPanicEvent      = "app:panic"
)

type BackendErrorPayload struct {
        Kind       string    `json:"kind"`
        Message    string    `json:"message"`
        Time       time.Time `json:"time"`
        StackTrace string    `json:"stackTrace,omitempty"`
}

func emit(name string, payload BackendErrorPayload) {
        app := application.Get()
        if app == nil {
                return
        }

        app.Event.Emit(name, payload)
}

// EmitFatalError reports a severe backend error without crashing the process.
func EmitFatalError(kind string, message string, stackTrace string) {
        emit(BackendFatalErrorEvent, BackendErrorPayload{
                Kind:       kind,
                Message:    message,
                Time:       time.Now(),
                StackTrace: stackTrace,
        })
}

// EmitRecoveredPanic reports a recovered backend panic to the frontend.
func EmitRecoveredPanic(details *application.PanicDetails) {
        if details == nil {
                return
        }

        emit(BackendPanicEvent, BackendErrorPayload{
                Kind:       "panic",
                Message:    fmt.Sprint(details.Error),
                Time:       details.Time,
                StackTrace: details.StackTrace,
        })
}
