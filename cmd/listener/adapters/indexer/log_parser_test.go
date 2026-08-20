package indexer

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
)

func TestNewLogParser(t *testing.T) {
	client := &ethclient.Client{}
	cfg := &config.Config{
		PrivateHub: config.PrivateHub{
			TokenCore:                  "0x11",
			EnygmaTokenManager:         "0x12",
			ParticipantCore:            "0x13",
			AuditManager:               "0x14",
			TokenRegistry:              "0x01",
			Teleport:                   "0x02",
			ParticipantStorageContract: "0x03",
			ProofsAddress:              "0x05",
			EnygmaPNHEvents:            "0x06",
			EnygmaTeleport:             "0x07",
			DvpTeleport:                "0x08",
			TokenFreezeManager:         "0x09",
		},
	}

	stubLogger := &testutil.StubLogger{}

	// Check if NewLogParser initializes correctly
	lp, err := NewLogParser(client, &Contracts{}, &testutil.StubProvider{}, cfg, stubLogger)
	if err != nil {
		t.Fatalf("NewLogParser returned error: %v", err)
	}
	if lp == nil {
		t.Fatalf("NewLogParser returned nil parser")
	}
	if lp.Client != client {
		t.Fatalf("Client mismatch: got %p, want %p", lp.Client, client)
	}
	if lp.Config != cfg {
		t.Fatalf("Config pointer mismatch")
	}
	if lp.Contracts == nil {
		t.Fatalf("Contracts not initialized")
	}
	if got, want := len(lp.AddressToContractParsers), 9; got != want {
		t.Fatalf("registry size = %d, want %d", got, want)
	}

	// Check if each expected contract is present in the registry
	expected := []struct{ addr, name string }{
		{cfg.PrivateHub.TokenCore, "TokenCore"},
		{cfg.PrivateHub.EnygmaTokenManager, "EnygmaTokenManager"},
		{cfg.PrivateHub.Teleport, "Teleport"},
		{cfg.PrivateHub.ParticipantCore, "ParticipantCore"},
		{cfg.PrivateHub.AuditManager, "AuditManager"},
		{cfg.PrivateHub.ProofsAddress, "Proofs"},
		{cfg.PrivateHub.EnygmaTeleport, "EnygmaTeleport"},
		{cfg.PrivateHub.DvpTeleport, "DvpTeleport"},
		{cfg.PrivateHub.TokenFreezeManager, "TokenFreezeManager"},
	}
	for _, e := range expected {
		key := common.HexToAddress(e.addr)
		contractParsers, ok := lp.AddressToContractParsers[key]
		if !ok {
			t.Fatalf("missing entry for address %s", e.addr)
		}
		if contractParsers.Name != e.name {
			t.Fatalf("name for %s = %q, want %q", e.addr, contractParsers.Name, e.name)
		}
	}

	// Check if each contract has EventParsers registered
	for addr, contractParsers := range lp.AddressToContractParsers {
		if len(contractParsers.Parsers) == 0 {
			t.Errorf("Contract %s (%s) has no EventParsers registered", contractParsers.Name, addr.String())
		} else {
			for _, ep := range contractParsers.Parsers {
				if ep.eventName == "" {
					t.Errorf("EventParser in %s (%s) has empty eventName", contractParsers.Name, addr.String())
				}
				if ep.eventSignatureHash == (common.Hash{}) {
					t.Errorf("EventParser in %s (%s) has empty eventSignature", contractParsers.Name, addr.String())
				}
				if ep.parser == nil {
					t.Errorf("EventParser in %s (%s) has nil parser function", contractParsers.Name, addr.String())
				}
			}
		}
	}
}

//nolint:gocognit // test function requires exhaustive case coverage; splitting would reduce readability
func TestLogParser_ParseLogs(t *testing.T) {
	client := &ethclient.Client{}
	cfg := &config.Config{PrivateHub: config.PrivateHub{
		TokenCore:          "0x11",
		EnygmaTokenManager: "0x12",
		Teleport:           "0x02",
		ParticipantCore:    "0x13",
		AuditManager:       "0x14",
		ProofsAddress:      "0x05",
		EnygmaTeleport:     "0x07",
	}}

	stubLogger := &testutil.StubLogger{}

	// Build registry manually after creating parser
	lp, err := NewLogParser(client, &Contracts{}, &testutil.StubProvider{Timestamp: 1000}, cfg, stubLogger)
	if err != nil {
		t.Fatalf("unexpected err creating log parser: %v", err)
	}

	// Override registry with a custom dummy parser to avoid nil contract method calls
	dummyAddr := common.HexToAddress("0x9000000000000000000000000000000000000001")
	dummySig := common.HexToHash("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	lp.AddressToContractParsers = map[common.Address]ContractParsers{
		dummyAddr: {Name: "Dummy", Parsers: []EventParser{{
			eventName:          "DummyEvent",
			eventSignatureHash: dummySig,
			parser: func(l ethTypes.Log) (any, error) {
				if l.TxHash == common.HexToHash("0x456") {
					return nil, errors.New("boom")
				}
				return struct{ OK bool }{OK: true}, nil
			},
		}}},
	}

	goodLog := ethTypes.Log{
		Address:     dummyAddr,
		Topics:      []common.Hash{dummySig},
		BlockNumber: 10,
		TxHash:      common.HexToHash("0xabc"),
	}
	malformedLog := ethTypes.Log{
		Address:     dummyAddr,
		Topics:      []common.Hash{},
		BlockNumber: 11,
		TxHash:      common.HexToHash("0xdef"),
	}
	unwatchedLog := ethTypes.Log{
		Address:     dummyAddr,
		Topics:      []common.Hash{common.HexToHash("0xdeadbeef")},
		BlockNumber: 12,
		TxHash:      common.HexToHash("0x123"),
	}
	parseErrLog := ethTypes.Log{
		Address:     dummyAddr,
		Topics:      []common.Hash{dummySig},
		BlockNumber: 13,
		TxHash:      common.HexToHash("0x456"),
	}

	tests := []struct {
		name      string
		logs      []ethTypes.Log
		wantCount int
	}{
		{
			name:      "all good",
			logs:      []ethTypes.Log{goodLog},
			wantCount: 1,
		},
		{
			name:      "malformed skipped",
			logs:      []ethTypes.Log{malformedLog},
			wantCount: 0,
		},
		{
			name:      "unwatched skipped",
			logs:      []ethTypes.Log{unwatchedLog},
			wantCount: 0,
		},
		{
			name:      "parser error skipped",
			logs:      []ethTypes.Log{parseErrLog},
			wantCount: 0,
		},
		{
			name:      "mix only good counted",
			logs:      []ethTypes.Log{malformedLog, unwatchedLog, parseErrLog, goodLog},
			wantCount: 1,
		},
	}

	// Test parseSingleLog logic indirectly by constructing logs and invoking parseSingleLog directly.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stubLogger := &testutil.StubLogger{}

			// Recreate parser with fresh stub for this subtest
			lp.log = stubLogger

			count := 0
			for _, lg := range tc.logs {
				parsed, err := lp.parseSingleLog(lg, time.Time{})
				if len(lg.Topics) == 0 { // malformed
					if err == nil {
						t.Fatalf("expected error for malformed log")
					}
					continue
				}
				if len(lg.Topics) > 0 && lg.Topics[0] != dummySig { // unwatched
					if err == nil {
						t.Fatalf("expected error for unwatched event")
					}
					continue
				}
				if lg.TxHash == parseErrLog.TxHash { // parser error
					if err == nil {
						t.Fatalf("expected parser error")
					}
					continue
				}
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				if parsed == nil {
					t.Fatalf("parsed log nil")
				}
				count++
			}
			if count != tc.wantCount {
				t.Fatalf("count=%d want=%d", count, tc.wantCount)
			}
		})
	}
}
