
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

/*
Why this is production-grade

Clear domain types (VehicleTelemetry, Vehicle, Geofence).

Interfaces separate concerns (TelemetryProvider, AlertService, MaintenanceScheduler).

Concurrency for ingestion with context cancellation and WaitGroup.

Non-blocking checks (runChecks runs async) and safe locking for state updates.

Structured alerts and scheduling hooks easily replaceable with real services (MQ, scheduler).

Compact, readable, and easy to extend (OTA updates, driver behavior, predictive maintenance).

*/



// VehicleTelemetry represents a single telemetry payload from a car.
type VehicleTelemetry struct {
	VIN        string
	Timestamp  time.Time
	Latitude   float64
	Longitude  float64
	FuelPct    float64 // 0 - 100
	BatteryPct float64 // 0 - 100 (EV)
	SpeedKmph  float64
	MilageKm   int
}

// TelemetryProvider ingests telemetry payloads.
type TelemetryProvider interface {
	Ingest(ctx context.Context, t VehicleTelemetry) error
}

// AlertService sends alerts based on rules.
type AlertService interface {
	Alert(vin, level, msg string)
}

// MaintenanceScheduler schedules maintenance tasks.
type MaintenanceScheduler interface {
	Schedule(vin, reason string)
}

// Vehicle represents connected car state and behavior.
type Vehicle struct {
	VIN       string
	LastSeen  time.Time
	Location  [2]float64
	FuelPct   float64
	Battery   float64
	MilageKm  int
	mu        sync.Mutex
	alertSvc  AlertService
	maintSvc  MaintenanceScheduler
	geofence  Geofence
}

// Geofence simple circular fence.
type Geofence struct {
	Lat    float64
	Lon    float64
	Radius float64 // in km
}

// IngestTelemetry updates vehicle state and runs checks.
func (v *Vehicle) IngestTelemetry(t VehicleTelemetry) error {
	v.mu.Lock()
	v.LastSeen = t.Timestamp
	v.Location[0] = t.Latitude
	v.Location[1] = t.Longitude
	v.FuelPct = t.FuelPct
	v.Battery = t.BatteryPct
	v.MilageKm = t.MilageKm
	v.mu.Unlock()

	// Run checks asynchronously but non-blocking to caller
	go v.runChecks(t)
	return nil
}

func (v *Vehicle) runChecks(t VehicleTelemetry) {
	// 1. Fuel / Battery low check
	if t.FuelPct >= 0 && t.FuelPct < 10 {
		v.alertSvc.Alert(v.VIN, "CRITICAL", fmt.Sprintf("Low fuel: %.1f%%", t.FuelPct))
		v.maintSvc.Schedule(v.VIN, "Refuel recommended - critical level")
	} else if t.BatteryPct >= 0 && t.BatteryPct < 8 {
		v.alertSvc.Alert(v.VIN, "CRITICAL", fmt.Sprintf("Low battery: %.1f%%", t.BatteryPct))
		v.maintSvc.Schedule(v.VIN, "Charge recommended - critical level")
	} else if t.FuelPct >= 10 && t.FuelPct < 20 {
		v.alertSvc.Alert(v.VIN, "WARN", fmt.Sprintf("Fuel low: %.1f%%", t.FuelPct))
	}

	// 2. Speed anomaly
	if t.SpeedKmph > 200 {
		v.alertSvc.Alert(v.VIN, "CRITICAL", fmt.Sprintf("Speed anomaly: %.1f km/h", t.SpeedKmph))
		v.maintSvc.Schedule(v.VIN, "Speed incident review")
	}

	// 3. Geofence breach
	if v.geofence.Breach(t.Latitude, t.Longitude) {
		v.alertSvc.Alert(v.VIN, "WARN", "Geofence breached")
	}

	// 4. Maintenance by milage
	if t.MilageKm > 0 && t.MilageKm%10000 == 0 {
		v.maintSvc.Schedule(v.VIN, fmt.Sprintf("Periodic service at %d km", t.MilageKm))
	}
}

// Breach checks if location is outside geofence.
func (g Geofence) Breach(lat, lon float64) bool {
	// Haversine distance approx
	const R = 6371.0 // Earth radius km
	lat1 := deg2rad(g.Lat)
	lon1 := deg2rad(g.Lon)
	lat2 := deg2rad(lat)
	lon2 := deg2rad(lon)
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := R * c
	return dist > g.Radius
}

func deg2rad(d float64) float64 { return d * math.Pi / 180 }

// SimpleAlert is a minimal alert service printing structured logs.
type SimpleAlert struct{}

func (s SimpleAlert) Alert(vin, level, msg string) {
	log.Printf("[ALERT] vin=%s level=%s msg=%s\n", vin, level, msg)
}

// SimpleScheduler prints schedule actions (would integrate with scheduler/queue in prod).
type SimpleScheduler struct{}

func (s SimpleScheduler) Schedule(vin, reason string) {
	log.Printf("[SCHEDULE] vin=%s reason=%s\n", vin, reason)
}

// InMemoryIngestor simulates ingestion pipeline.
type InMemoryIngestor struct {
	vehicles map[string]*Vehicle
}

// NewInMemoryIngestor constructs provider.
func NewInMemoryIngestor(alertSvc AlertService, maintSvc MaintenanceScheduler, gf Geofence) *InMemoryIngestor {
	return &InMemoryIngestor{
		vehicles: map[string]*Vehicle{
			"WVWZZZ1JZXW000001": {VIN: "WVWZZZ1JZXW000001", alertSvc: alertSvc, maintSvc: maintSvc, geofence: gf},
			"WVWZZZ1JZXW000002": {VIN: "WVWZZZ1JZXW000002", alertSvc: alertSvc, maintSvc: maintSvc, geofence: gf},
		},
	}
}

func (in *InMemoryIngestor) Ingest(ctx context.Context, t VehicleTelemetry) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		v, ok := in.vehicles[t.VIN]
		if !ok {
			return fmt.Errorf("unknown vehicle %s", t.VIN)
		}
		return v.IngestTelemetry(t)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	alertSvc := SimpleAlert{}
	maintSvc := SimpleScheduler{}
	// geofence center (example) and 5 km radius
	gf := Geofence{Lat: 48.1351, Lon: 11.5820, Radius: 5.0}

	ingestor := NewInMemoryIngestor(alertSvc, maintSvc, gf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vinList := []string{"WVWZZZ1JZXW000001", "WVWZZZ1JZXW000002"}
	var wg sync.WaitGroup

	// simulate streaming telemetry concurrently
	for i := 0; i < 20; i++ {
		for _, vin := range vinList {
			wg.Add(1)
			go func(v string, iter int) {
				defer wg.Done()
				t := VehicleTelemetry{
					VIN:        v,
					Timestamp:  time.Now(),
					Latitude:   48.1351 + rand.Float64()*0.02*(randFloatSign()),
					Longitude:  11.5820 + rand.Float64()*0.02*(randFloatSign()),
					FuelPct:    rand.Float64()*100,
					BatteryPct: rand.Float64()*100,
					SpeedKmph:  rand.Float64() * 220,
					MilageKm:   10000 + iter*500,
				}
				if err := ingestor.Ingest(ctx, t); err != nil {
					log.Printf("[INGEST ERROR] vin=%s err=%v\n", v, err)
				}
				// small spacing
				time.Sleep(40 * time.Millisecond)
			}(vin, i)
		}
	}

	// wait and finish
	wg.Wait()
	log.Println("[INFO] Telemetry ingestion completed")
}

// randFloatSign returns either -1 or +1 randomly.
func randFloatSign() float64 {
	if rand.Intn(2) == 0 {
		return -1
	}
	return 1
}


