package destinationharness

//go:generate go run github.com/synapsecns/sanguine/tools/abigen generate --sol  ../../../../packages/contracts/flattened/DestinationHarness.sol --pkg destinationharness --sol-version 0.8.13 --filename destinationharness
// line after go:generate cannot be left blank