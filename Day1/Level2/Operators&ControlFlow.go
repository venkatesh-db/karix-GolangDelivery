package main

/*

Use if for fast failure and switch for environment/config logic.

Loops with counters or ranges for repeated work.

Logging prefixes [ALERT], [WARN], [OK] make logs parseable in monitoring systems.

*/


import (
    "fmt"
    "log"
    "time"
)

func main() {
    serviceHealthy := true
    requestLatency := 82 // in milliseconds

    // Decision control
    if !serviceHealthy {
        log.Println("[ALERT] Service is unhealthy — taking corrective action.")
    } else if requestLatency > 100 {
        log.Printf("[WARN] High latency detected: %dms", requestLatency)
    } else {
        log.Printf("[OK] Service stable. Latency=%dms", requestLatency)
    }

    // Switch-based environment configuration
    env := "staging"

    switch env {
    case "dev":
        log.Println("[INIT] Loading Dev Configuration...")
    case "staging":
        log.Println("[INIT] Loading Staging Configuration...")
    case "prod":
        log.Println("[INIT] Loading Production Configuration...")
    default:
        log.Println("[INIT] Unknown Environment — Exiting.")
    }

    // For loop to simulate request processing
    for i := 1; i <= 5; i++ {
        fmt.Printf("Processing request #%d\n", i)
        time.Sleep(200 * time.Millisecond)
    }
}
