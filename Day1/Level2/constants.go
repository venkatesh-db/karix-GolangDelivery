
package main

import "log"



/*

🟢 Production Conventions

Use constants for static configurations.

Use iota for environment or status enums.

Use type aliases (Milliseconds) for strong clarity and unit safety.

Consistent camelCase and [CONFIG] log tags.

*/

// Constant configuration (never changes at runtime)
const (
    MaxRetries      = 3
    DefaultTimeoutS = 5
)

// Enum for environment
type Environment int

const (
    EnvDev Environment = iota
    EnvStage
    EnvProd
)

// Type alias for clarity
type Milliseconds int64

func main() {
    var env Environment = EnvProd
    var apiTimeout Milliseconds = 1500

    log.Printf("[CONFIG] Environment=%d | Timeout=%dms | Retries=%d",
        env, apiTimeout, MaxRetries)
}
