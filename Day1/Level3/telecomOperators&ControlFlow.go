package main

import (
    "log"
)

/*

🟢 Telecom Insights

Decision rules help with signal optimization.

Switches manage status enums.

Loops simulate multi-sector scanning.

*/

func main() {
    signalStrength := -78.5 // dBm
    towerStatus := "ACTIVE"

    // Decision making
    if signalStrength < -85.0 {
        log.Printf("[TELECOM] Weak Signal: %.1fdBm | Action: Boost Power", signalStrength)
    } else if signalStrength >= -85.0 && signalStrength < -70.0 {
        log.Printf("[TELECOM] Moderate Signal: %.1fdBm | Action: Maintain Stability", signalStrength)
    } else {
        log.Printf("[TELECOM] Strong Signal: %.1fdBm | Action: Optimize Load", signalStrength)
    }

    // Switch tower operational mode
    switch towerStatus {
    case "ACTIVE":
        log.Println("[TELECOM] Tower operating normally.")
    case "MAINTENANCE":
        log.Println("[TELECOM] Tower under scheduled maintenance.")
    default:
        log.Println("[TELECOM] Tower status unknown — investigating.")
    }

    // Loop through sectors to monitor
    for sector := 1; sector <= 3; sector++ {
        log.Printf("[TELECOM] Scanning Sector-%d — OK", sector)
    }
}
