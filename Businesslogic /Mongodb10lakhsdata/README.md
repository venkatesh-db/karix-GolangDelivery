# NPCI-Style MongoDB Seeder

A Go-based simulator that floods MongoDB with ~2.5 million production-grade payment transactions (25 lakh) resembling NPCI/UPI/IMPS traffic. The generator deliberately injects common production issues such as duplicate references, SLA breaches, stale statuses, AML hits, and device compromises so downstream systems can be validated under stress.

## Features

- **Configurable scale**: defaults to 2.5 million docs with tunable workers, batch sizes, and ratios via env vars.
- **Rich model**: payer/payee, device fingerprint, settlement windows, compliance flags, routing metadata, failure semantics, etc.
- **Issue simulation**: duplicates, stale pending items, velocity/AML alerts, high-value spikes, switch failures.
- **Concurrent writers**: worker pool bulk inserts using unordered batches for throughput similar to production traffic.
- **Index bootstrap**: auto-creates practical indexes (txn_id, utr, status, AML flags) for realistic query planners.

## Quick Start

1. **Prerequisites**
   - Go 1.21+
   - MongoDB 6.0+ reachable at the configured URI

2. **Environment (optional overrides)**

   | Variable | Default | Description |
   |----------|---------|-------------|
   | `MONGO_URI` | `mongodb://localhost:27017` | Target MongoDB connection string |
   | `MONGO_DB` | `npcisim` | Database for seeded data |
   | `MONGO_COLLECTION` | `transactions` | Collection name |
   | `TOTAL_RECORDS` | `2500000` | Number of documents to create |
   | `BATCH_SIZE` | `2000` | InsertMany batch size |
   | `WORKERS` | `8` | Concurrent generator goroutines |
   | `FAILURE_RATIO` | `0.12` | Probability of FAILED/TIMEOUT statuses |
   | `COMPLIANCE_RATIO` | `0.03` | Ratio of AML/compliance hits |
   | `DUPLICATE_RATIO` | `0.02` | Ratio of duplicate reference IDs |
   | `STALE_STATUS_RATIO` | `0.04` | Ratio of stuck `PENDING` items moved back in time |
   | `LOGGER_INTERVAL` | `5s` | Progress log cadence |

3. **Install dependencies & run**

   ```bash
   cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /Mongodb10lakhsdata"
   go mod tidy
   go run ./cmd/seed
   ```

4. **Start MongoDB (if needed)**

   ```bash
   brew services start mongodb-community@7.0
   # or run your mongod --dbpath ...
   ```

5. **Observed output** (local MacBook Pro + MongoDB 7.0)

   ```text
   2025/11/29 14:02:19 starting generator: total=2500000 workers=8 batch=2000
   2025/11/29 14:02:24 progress: 510000/2500000 (20.40%)
   2025/11/29 14:02:29 progress: 1008000/2500000 (40.32%)
   2025/11/29 14:02:34 progress: 1480000/2500000 (59.20%)
   2025/11/29 14:02:39 progress: 1940000/2500000 (77.60%)
   2025/11/29 14:02:44 progress: 2406000/2500000 (96.24%)
   2025/11/29 14:02:45 finished inserting 2500000 docs (failed:0) in 25.98s
   ```

## Data Shape

Each `transactions` document roughly looks like:

```json
{
  "txn_id": "NPCI0000000123456",
  "reference_id": "REF0000000123456",
  "utr": "HDFC29011 2345678",
  "payment_rail": "UPI",
  "instrument_type": "VPA",
  "channel": "P2M",
  "amount": { "value_paise": 755000, "currency": "INR" },
  "payer": { "customer_id": "CUST00001234", "bank_ifsc": "HDFC0000123", ... },
  "payee": { ... },
  "status": "FAILED",
  "failure": { "code": "U003", "category": "NETWORK" },
  "settlement": { "window": "T+0 15:00", "recon_status": "NOT_REQUIRED" },
  "device": { "device_type": "Android", "is_compromised": true },
  "compliance_flags": { "aml_hit": true, "list_match": "UNCFT" },
  "anomalies": ["compliance_hit", "stale_status"],
  "metadata": { "issuer_code": "HDFC", "switch": "NPCI_CORE" },
  "created_at": ISODate(...),
  "updated_at": ISODate(...)
}
```

This mirrors real monitoring scenarios: AML teams can query `compliance_flags.aml_hit`, recon squads can slice by `settlement.recon_status`, and payment reliability teams can replay failure modes.

## Tuning ideas

- Lower `TOTAL_RECORDS` plus `BATCH_SIZE` for smoke tests.
- Increase `WORKERS` after observing Mongo server CPU/IO headroom.
- Crank `FAILURE_RATIO` or `COMPLIANCE_RATIO` when stress-testing negative workflows.
- Point to a Mongo Atlas sharded cluster to test chunk migrations or balancer behavior.

## Cleanup

Drop the collection when done:

```bash
mongosh "$MONGO_URI" --eval 'db.getSiblingDB("npcisim").transactions.drop()'
```

Happy load-testing!
