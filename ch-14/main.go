package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// 1
func middlewareWithTimeout(timeout int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Millisecond)
			defer cancel()
			r = r.WithContext(newCtx)
			h.ServeHTTP(w, r)
		})
	}
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte("Accepted"))
	case <-r.Context().Done():
		w.Write([]byte("Timeout"))
	}
}

// 2
func generateNums(ctx context.Context) (sum int, iter int, reason string) {
outer:
	for {
		ranNumOne, ranNumTwo := rand.Intn(100000000), rand.Intn(100000000)
		sum = ranNumOne + ranNumTwo
		iter++
		select {
		case <-ctx.Done():
			reason = "Timeout"
			break outer
		default:
			if sum == 1234 {
				reason = "Sum equals to 1234"
				break outer
			}
		}
	}

	return sum, iter, reason
}

// 3
type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
)

func writeJLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, "log_level", level)
}

func getJLevel(ctx context.Context) Level {
	val, ok := ctx.Value("log_level").(Level)
	if !ok {
		var retLevel Level
		return retLevel
	}
	return val
}

func logMiddleware(h http.Handler, logLevel Level) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level := r.URL.Query().Get("log_level")
		r = r.WithContext(writeJLevel(r.Context(), logLevel))
		Log(r.Context(), Level(level), "Received request")
		h.ServeHTTP(w, r)
	})
}

func Log(ctx context.Context, level Level, message string) {
	var inLevel = getJLevel(ctx)

	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}
	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}

func main() { // 2
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s, i, r := generateNums(timeoutCtx)
	fmt.Println(s, i, r)
	// 1 + 3
	timeoutMw := middlewareWithTimeout(rand.Intn(100000000))
	mux := http.NewServeMux()
	mux.Handle("/", logMiddleware(timeoutMw(http.HandlerFunc(randomHandler)), Debug))
	http.ListenAndServe("localhost:8080", mux)
}
