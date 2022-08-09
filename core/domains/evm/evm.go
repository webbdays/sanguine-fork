// Package evm TODO description
package evm

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/synapsecns/sanguine/core/config"
	"github.com/synapsecns/sanguine/core/domains"
	"github.com/synapsecns/synapse-node/pkg/evm"
)

type evmClient struct {
	// name is the name of the evm client
	name string
	// config is the config of the evm client
	config config.DomainConfig
	// client uses the old synapse client for now
	client evm.Chain
	// home contains the home contract
	home domains.HomeContract
	// attestationCollecotr contains the attestation collector contract
	attestationCollector domains.AttestationCollectorContract
}

var _ domains.DomainClient = &evmClient{}

// NewEVM creates a new evm client.
func NewEVM(ctx context.Context, name string, domain config.DomainConfig) (domains.DomainClient, error) {
	underlyingClient, err := evm.NewFromURL(ctx, domain.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("could not get evm: %w", err)
	}

	boundHome, err := NewHomeContract(ctx, underlyingClient, common.HexToAddress(domain.HomeAddress))
	if err != nil {
		return nil, fmt.Errorf("could not bind home contract: %w", err)
	}

	boundCollector, err := NewAttestationCollectorContract(ctx, underlyingClient, common.HexToAddress(domain.AttesationCollectorAddress))
	if err != nil {
		return nil, fmt.Errorf("could not bind attestation contract: %w", err)
	}

	return evmClient{
		name:                 name,
		config:               domain,
		client:               underlyingClient,
		attestationCollector: boundCollector,
		home:                 boundHome,
	}, nil
}

// Name gets the name of the evm client.
func (e evmClient) Name() string {
	return e.name
}

// Config gets the config the evm client was initiated with.
func (e evmClient) Config() config.DomainConfig {
	return e.config
}

// BlockNumber gets the latest block number.
func (e evmClient) BlockNumber(ctx context.Context) (uint32, error) {
	blockNumber, err := e.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not get block number: %w", err)
	}

	return uint32(blockNumber), nil
}

// Home returns the bound home contract.
func (e evmClient) Home() domains.HomeContract {
	return e.home
}

// AttestationCollector gets the attestation collector.
func (e evmClient) AttestationCollector() domains.AttestationCollectorContract {
	return e.attestationCollector
}