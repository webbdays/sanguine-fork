// Package manager manages deployers to make them as simple as possible
package manager

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/synapsecns/sanguine/ethergo/deployer"
	"github.com/synapsecns/synapse-node/testutils/backends"
	"sync"
	"testing"
)

// DeployerManager is responsible for wrapping contract registry with easy to use getters that correctly cast the handles.
// since ContractRegistry is meant to be kept pure and go does not support generics, the sole function is to provide
// handler wrappers around the registry. This will no longer be required when go supports generics: https://blog.golang.org/generics-proposal
type DeployerManager struct {
	// t is the testing object
	t *testing.T
	// registries stores the contract registries
	registries map[string]deployer.ContractRegistry
	// structMux prevents race conditions
	structMux sync.RWMutex
	// deployers adds a list of default deployers
	deployers []DeployerFunc
}

// DeployerFunc defines a deployer we can use.
type DeployerFunc func(registry deployer.GetOnlyContractRegistry, backend backends.SimulatedTestBackend) deployer.ContractDeployer

// NewDeployerManager creates a new deployment helper.
func NewDeployerManager(t *testing.T, deployers ...DeployerFunc) (d *DeployerManager) {
	t.Helper()
	d = &DeployerManager{
		t:          t,
		structMux:  sync.RWMutex{},
		registries: make(map[string]deployer.ContractRegistry),
		deployers:  deployers,
	}
	return d
}

// T is the testing object.
func (d *DeployerManager) T() *testing.T {
	return d.t
}

// SetT sets the testing object.
func (d *DeployerManager) SetT(t *testing.T) {
	t.Helper()
	d.t = t
}

// BulkDeploy synchronously deploys a bunch of contracts as quickly as possible to speed up tests.
// in a future version this will utilize dependency trees. Returns nothing when complete.
func (d *DeployerManager) BulkDeploy(ctx context.Context, testBackends []backends.SimulatedTestBackend, contracts ...deployer.ContractType) {
	wg := sync.WaitGroup{}
	for _, backend := range testBackends {
		wg.Add(1)
		go func(backend backends.SimulatedTestBackend) {
			defer wg.Done()
			cr := d.GetContractRegistry(backend)

			for _, contract := range contracts {
				cr.Get(ctx, contract)
			}
		}(backend)
	}
	wg.Wait()
}

// GetContractRegistry gets a contract registry for a backend and creates it if it does not exist.
func (d *DeployerManager) GetContractRegistry(backend backends.SimulatedTestBackend) deployer.ContractRegistry {
	d.structMux.Lock()
	defer d.structMux.Unlock()
	// if registry exists, return it
	contractRegistry, ok := d.registries[backend.GetBigChainID().String()]
	if ok {
		return contractRegistry
	}

	contractRegistry = deployer.NewContractRegistry(d.t, backend)

	for _, d := range d.deployers {
		contractRegistry.RegisterContractDeployer(d(contractRegistry, backend))
	}

	d.registries[backend.GetBigChainID().String()] = contractRegistry
	return contractRegistry
}

// Get gets the contract from the registry.
func (d *DeployerManager) Get(ctx context.Context, backend backends.SimulatedTestBackend, contractType deployer.ContractType) backends.DeployedContract {
	return d.GetContractRegistry(backend).Get(ctx, contractType)
}

var _ suite.TestingSuite = &DeployerManager{}