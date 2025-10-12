package middleware

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestLoggerMiddleware(t *testing.T) {
    var buf bytes.Buffer
    old := log.Default().Writer()
    log.SetOutput(&buf)
    defer log.SetOutput(old)

    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, "ok")
    })

    rr := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/test", nil)
    Logger(next).ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
    s := buf.String()
    if !bytes.Contains([]byte(s), []byte("Started GET /test")) {
        t.Fatalf("expected start log, got: %s", s)
    }
    if !bytes.Contains([]byte(s), []byte("Completed GET /test")) {
        t.Fatalf("expected completion log, got: %s", s)
    }
}

func TestRecovererMiddleware(t *testing.T) {
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        panic("boom")
    })

    rr := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/panic", nil)
    Recoverer(next).ServeHTTP(rr, req)

    if rr.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", rr.Code)
    }
}

func TestInjectActorUser(t *testing.T) {
    actor := "user-123"
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id, ok := ActorUserIDFromContext(r.Context())
        if !ok || id != actor {
            t.Fatalf("expected actor %s, got %s (ok=%v)", actor, id, ok)
        }
        w.WriteHeader(http.StatusOK)
    })

    rr := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/actor", nil)
    InjectActorUser(actor)(next).ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
}