package generator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/venkatesh/mongodb-simulator/internal/config"
	"github.com/venkatesh/mongodb-simulator/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Result summarises a generator run.
type Result struct {
	Inserted int64
	Failed   int64
	Duration time.Duration
}

// Generator owns fabrication logic + Mongo persistence.
type Generator struct {
	cfg          *config.Config
	collection   *mongo.Collection
	seq          atomic.Int64
	inserted     atomic.Int64
	failed       atomic.Int64
	refStore     []string
	refMu        sync.RWMutex
	banks        []bank
	paymentRails []string
	instrument   []string
	channels     []string
	deviceTypes  []string
	appVersions  []string
	failureCat   []models.FailureDetail
	settleWindow []string
	firstNames   []string
	lastNames    []string
	kycStates    []string
}

type bank struct {
	Code string
	IFSC string
}

type batchJob struct {
	size int
}

// New constructs a ready-to-run generator instance.
func New(cfg *config.Config, coll *mongo.Collection) *Generator {
	return &Generator{
		cfg:        cfg,
		collection: coll,
		banks: []bank{
			{Code: "HDFC", IFSC: "HDFC0000123"},
			{Code: "ICIC", IFSC: "ICIC0000456"},
			{Code: "SBI", IFSC: "SBIN0000800"},
			{Code: "PNB", IFSC: "PUNB0003300"},
			{Code: "AXIS", IFSC: "UTIB0000505"},
			{Code: "YESB", IFSC: "YESB0000999"},
			{Code: "UBIN", IFSC: "UBIN0532154"},
		},
		paymentRails: []string{"UPI", "IMPS", "NEFT", "RTGS", "AEPS"},
		instrument:   []string{"VPA", "ACCOUNT", "AADHAAR", "CARD", "NETC"},
		channels:     []string{"P2P", "P2M", "AutoPay", "Mandate", "IPO", "BillPay"},
		deviceTypes:  []string{"Android", "iOS", "POS", "MicroATM"},
		appVersions:  []string{"4.18.2", "5.0.0", "4.19.7", "5.1.3", "6.0.1"},
		failureCat: []models.FailureDetail{
			{Code: "U003", Category: "NETWORK", Severity: "HIGH", Description: "Switch timeout"},
			{Code: "U017", Category: "REVERSAL", Severity: "MEDIUM", Description: "Issuer pending"},
			{Code: "KYC12", Category: "KYC", Severity: "LOW", Description: "KYC expired"},
			{Code: "LIM45", Category: "RISK", Severity: "HIGH", Description: "Velocity breached"},
		},
		settleWindow: []string{"T+0 15:00", "T+0 20:00", "T+1 10:00"},
		firstNames:   []string{"Aarav", "Vihaan", "Ananya", "Advika", "Zara", "Kabir", "Reyansh", "Ira", "Rohan", "Myra"},
		lastNames:    []string{"Sharma", "Patel", "Reddy", "Iyer", "Khan", "Singh", "Das", "Varma", "Nayak", "Shetty"},
		kycStates:    []string{"FULL", "MINIMAL", "REVOKED"},
	}
}

// Run performs the configured number of inserts.
func (g *Generator) Run(ctx context.Context) (*Result, error) {
	start := time.Now()
	jobs := make(chan batchJob)
	errCh := make(chan error, 1)
	progressCtx, cancelProgress := context.WithCancel(ctx)
	defer cancelProgress()

	var wg sync.WaitGroup
	for workerID := 0; workerID < g.cfg.Workers; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)*1_000_003))
			for job := range jobs {
				if err := g.handleBatch(progressCtx, r, job.size); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			}
		}(workerID)
	}

	go g.logProgress(progressCtx)

	total := g.cfg.TotalRecords
	for created := 0; created < total; created += g.cfg.BatchSize {
		remaining := total - created
		batchSize := g.cfg.BatchSize
		if remaining < batchSize {
			batchSize = remaining
		}

		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return nil, ctx.Err()
		case jobs <- batchJob{size: batchSize}:
		}
	}

	close(jobs)
	wg.Wait()
	cancelProgress()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return &Result{
		Inserted: g.inserted.Load(),
		Failed:   g.failed.Load(),
		Duration: time.Since(start),
	}, nil
}

func (g *Generator) handleBatch(parent context.Context, r *rand.Rand, batchSize int) error {
	docs := make([]interface{}, batchSize)
	for i := 0; i < batchSize; i++ {
		docs[i] = g.newTransaction(r)
	}

	opCtx, cancel := context.WithTimeout(parent, g.cfg.OperationTimeout)
	defer cancel()

	_, err := g.collection.InsertMany(opCtx, docs, options.InsertMany().SetOrdered(false))
	if err != nil {
		g.failed.Add(int64(batchSize))
		return fmt.Errorf("insert batch: %w", err)
	}

	g.inserted.Add(int64(batchSize))
	return nil
}

func (g *Generator) newTransaction(r *rand.Rand) models.Transaction {
	seq := g.seq.Add(1)
	issuer := g.banks[r.Intn(len(g.banks))]
	acquirer := g.banks[r.Intn(len(g.banks))]
	amountPaise := int64(10_000 + r.Intn(9_900_000)) // ₹100 to ₹99,000
	created := randomPastTime(r)
	updated := created.Add(time.Duration(r.Intn(900)) * time.Second)

	status := g.pickStatus(r)
	failure := models.FailureDetail{}
	statusReason := ""
	var completedAt *time.Time

	if status == "SUCCESS" || status == "REVERSED" {
		complete := updated.Add(time.Duration(r.Intn(120)) * time.Second)
		completedAt = &complete
	} else if status == "FAILED" || status == "TIMEOUT" {
		failure = g.failureCat[r.Intn(len(g.failureCat))]
		statusReason = failure.Description
	}

	referenceID, duplicate := g.referenceID(r, seq)
	utr := g.utr(r, issuer)

	complianceHit := r.Float64() < g.cfg.ComplianceRatio
	flags := models.ComplianceFlags{}
	anomalies := make([]string, 0, 4)
	if complianceHit {
		flags.AMLHit = r.Float64() < 0.7
		flags.GeoMismatch = r.Float64() < 0.4
		flags.VelocitySpike = r.Float64() < 0.5
		if flags.AMLHit {
			flags.ListMatch = "UNCFT"
		}
		anomalies = append(anomalies, "compliance_hit")
	}

	if duplicate {
		anomalies = append(anomalies, "duplicate_reference")
	}

	if amountPaise > 5_000_000 {
		anomalies = append(anomalies, "high_value")
	}

	if r.Float64() < g.cfg.StaleStatusRatio && status == "PENDING" {
		created = created.Add(-time.Duration(48+r.Intn(72)) * time.Hour)
		anomalies = append(anomalies, "stale_status")
	}

	payer := g.participant(r)
	payee := g.participant(r)

	device := g.device(r, complianceHit)
	settlement := g.settlement(r, created, amountPaise, status)

	retryCount := r.Intn(2)
	if status != "SUCCESS" {
		retryCount = 1 + r.Intn(3)
	}

	metadata := map[string]string{
		"issuer_code":   issuer.Code,
		"acquirer_code": acquirer.Code,
		"switch":        "NPCI_CORE",
		"route":         g.channels[r.Intn(len(g.channels))],
		"cluster":       fmt.Sprintf("%02d", r.Intn(24)),
	}

	return models.Transaction{
		TxnID:          fmt.Sprintf("NPCI%015d", seq),
		ReferenceID:    referenceID,
		UTR:            utr,
		PaymentRail:    g.paymentRails[r.Intn(len(g.paymentRails))],
		InstrumentType: g.instrument[r.Intn(len(g.instrument))],
		Channel:        metadata["route"],
		Amount: models.MonetaryAmount{
			ValuePaise: amountPaise,
			Currency:   "INR",
		},
		Payer:           payer,
		Payee:           payee,
		Status:          status,
		StatusReason:    statusReason,
		RetryCount:      retryCount,
		Failure:         failure,
		Settlement:      settlement,
		Device:          device,
		ComplianceFlags: flags,
		Anomalies:       anomalies,
		Metadata:        metadata,
		CreatedAt:       created,
		UpdatedAt:       updated,
		CompletedAt:     completedAt,
	}
}

func (g *Generator) participant(r *rand.Rand) models.Participant {
	first := g.firstNames[r.Intn(len(g.firstNames))]
	last := g.lastNames[r.Intn(len(g.lastNames))]
	bank := g.banks[r.Intn(len(g.banks))]
	fullName := fmt.Sprintf("%s %s", first, last)
	return models.Participant{
		CustomerID: fmt.Sprintf("CUST%08d", r.Intn(10_000_000)),
		Name:       fullName,
		BankIFSC:   bank.IFSC,
		Account:    fmt.Sprintf("%011d", r.Int63()%99999999999),
		VPA:        strings.ToLower(first) + ".pay@" + strings.ToLower(bank.Code),
		Latitude:   8 + r.Float64()*20,
		Longitude:  70 + r.Float64()*15,
		RiskScore:  r.Intn(100),
		KYCStatus:  g.kycStates[r.Intn(len(g.kycStates))],
	}
}

func (g *Generator) device(r *rand.Rand, compromised bool) models.DeviceProfile {
	return models.DeviceProfile{
		DeviceID:      fmt.Sprintf("DEV-%x", primitive.NewObjectID()),
		IPAddress:     fmt.Sprintf("%d.%d.%d.%d", r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255)),
		GeoHash:       fmt.Sprintf("%d%d%d", r.Intn(9), r.Intn(9), r.Intn(9)),
		DeviceType:    g.deviceTypes[r.Intn(len(g.deviceTypes))],
		AppVersion:    g.appVersions[r.Intn(len(g.appVersions))],
		OSVersion:     fmt.Sprintf("%d.%d.%d", 10+r.Intn(3), r.Intn(5), r.Intn(9)),
		IsCompromised: compromised && r.Float64() < 0.5,
	}
}

func (g *Generator) settlement(r *rand.Rand, created time.Time, amountPaise int64, status string) models.Settlement {
	window := g.settleWindow[r.Intn(len(g.settleWindow))]
	settleDate := created.Add(time.Duration(12+r.Intn(24)) * time.Hour)
	valueDate := settleDate.Add(24 * time.Hour)
	recon := "NOT_REQUIRED"
	if status == "SUCCESS" {
		recon = "RECONCILED"
	} else if status == "REVERSED" {
		recon = "REVERSAL_PENDING"
	}
	return models.Settlement{
		Window:          window,
		SettlementDate:  settleDate,
		ValueDate:       valueDate,
		NetSettlementRs: float64(amountPaise) / 100,
		ReconStatus:     recon,
	}
}

func (g *Generator) referenceID(r *rand.Rand, seq int64) (string, bool) {
	if r.Float64() < g.cfg.DuplicateRatio {
		if ref, ok := g.refFromPool(r); ok {
			return ref, true
		}
	}
	ref := fmt.Sprintf("REF%015d", seq)
	g.saveReference(ref)
	return ref, false
}

func (g *Generator) saveReference(ref string) {
	g.refMu.Lock()
	defer g.refMu.Unlock()
	const maxSize = 5000
	if len(g.refStore) >= maxSize {
		copy(g.refStore, g.refStore[1:])
		g.refStore[maxSize-1] = ref
		return
	}
	g.refStore = append(g.refStore, ref)
}

func (g *Generator) refFromPool(r *rand.Rand) (string, bool) {
	g.refMu.RLock()
	defer g.refMu.RUnlock()
	if len(g.refStore) == 0 {
		return "", false
	}
	return g.refStore[r.Intn(len(g.refStore))], true
}

func (g *Generator) utr(r *rand.Rand, issuer bank) string {
	return fmt.Sprintf("%s%s%07d", issuer.Code, time.Now().Format("0201"), r.Intn(9_999_999))
}

func (g *Generator) pickStatus(r *rand.Rand) string {
	roll := r.Float64()
	switch {
	case roll < 0.74:
		return "SUCCESS"
	case roll < 0.82:
		return "PENDING"
	case roll < 0.90:
		return "FAILED"
	case roll < 0.97:
		return "REVERSED"
	default:
		return "TIMEOUT"
	}
}

func (g *Generator) logProgress(ctx context.Context) {
	ticker := time.NewTicker(g.cfg.LoggerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			inserted := g.inserted.Load()
			pct := float64(inserted) / float64(g.cfg.TotalRecords) * 100
			log.Printf("progress: %d/%d (%.2f%%)", inserted, g.cfg.TotalRecords, pct)
		}
	}
}

func randomPastTime(r *rand.Rand) time.Time {
	daysBack := r.Intn(5)
	seconds := r.Intn(86_400)
	return time.Now().Add(-time.Duration(daysBack)*24*time.Hour - time.Duration(seconds)*time.Second)
}
