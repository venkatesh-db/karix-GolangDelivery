package models

import "time"

// Transaction captures an end-to-end payment flow resembling production NPCI traffic.
type Transaction struct {
	TxnID           string            `bson:"txn_id" json:"txn_id"`
	ReferenceID     string            `bson:"reference_id" json:"reference_id"`
	UTR             string            `bson:"utr" json:"utr"`
	PaymentRail     string            `bson:"payment_rail" json:"payment_rail"`
	InstrumentType  string            `bson:"instrument_type" json:"instrument_type"`
	Channel         string            `bson:"channel" json:"channel"`
	Amount          MonetaryAmount    `bson:"amount" json:"amount"`
	Payer           Participant       `bson:"payer" json:"payer"`
	Payee           Participant       `bson:"payee" json:"payee"`
	Status          string            `bson:"status" json:"status"`
	StatusReason    string            `bson:"status_reason,omitempty" json:"status_reason,omitempty"`
	RetryCount      int               `bson:"retry_count" json:"retry_count"`
	Failure         FailureDetail     `bson:"failure" json:"failure"`
	Settlement      Settlement        `bson:"settlement" json:"settlement"`
	Device          DeviceProfile     `bson:"device" json:"device"`
	ComplianceFlags ComplianceFlags   `bson:"compliance_flags" json:"compliance_flags"`
	Anomalies       []string          `bson:"anomalies" json:"anomalies"`
	Metadata        map[string]string `bson:"metadata" json:"metadata"`
	CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
	CompletedAt     *time.Time        `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

// MonetaryAmount keeps currency + paisa precision without floats.
type MonetaryAmount struct {
	ValuePaise int64  `bson:"value_paise" json:"value_paise"`
	Currency   string `bson:"currency" json:"currency"`
}

// Participant captures payer/payee attributes.
type Participant struct {
	CustomerID string  `bson:"customer_id" json:"customer_id"`
	Name       string  `bson:"name" json:"name"`
	BankIFSC   string  `bson:"bank_ifsc" json:"bank_ifsc"`
	Account    string  `bson:"account" json:"account"`
	VPA        string  `bson:"vpa" json:"vpa"`
	Latitude   float64 `bson:"latitude" json:"latitude"`
	Longitude  float64 `bson:"longitude" json:"longitude"`
	RiskScore  int     `bson:"risk_score" json:"risk_score"`
	KYCStatus  string  `bson:"kyc_status" json:"kyc_status"`
}

// FailureDetail mirrors ISO8583/UPI failure semantics.
type FailureDetail struct {
	Code        string `bson:"code" json:"code"`
	Category    string `bson:"category" json:"category"`
	Severity    string `bson:"severity" json:"severity"`
	Description string `bson:"description" json:"description"`
}

// Settlement covers recon + SLA windows.
type Settlement struct {
	ReconStatus     string    `bson:"recon_status" json:"recon_status"`
	Window          string    `bson:"window" json:"window"`
	SettlementDate  time.Time `bson:"settlement_date" json:"settlement_date"`
	ValueDate       time.Time `bson:"value_date" json:"value_date"`
	NetSettlementRs float64   `bson:"net_settlement_rs" json:"net_settlement_rs"`
}

// DeviceProfile introduces fraud vectors.
type DeviceProfile struct {
	DeviceID      string `bson:"device_id" json:"device_id"`
	IPAddress     string `bson:"ip_address" json:"ip_address"`
	GeoHash       string `bson:"geo_hash" json:"geo_hash"`
	DeviceType    string `bson:"device_type" json:"device_type"`
	AppVersion    string `bson:"app_version" json:"app_version"`
	OSVersion     string `bson:"os_version" json:"os_version"`
	IsCompromised bool   `bson:"is_compromised" json:"is_compromised"`
}

// ComplianceFlags highlights AML-style alerts.
type ComplianceFlags struct {
	AMLHit        bool   `bson:"aml_hit" json:"aml_hit"`
	GeoMismatch   bool   `bson:"geo_mismatch" json:"geo_mismatch"`
	VelocitySpike bool   `bson:"velocity_spike" json:"velocity_spike"`
	ListMatch     string `bson:"list_match" json:"list_match"`
}
