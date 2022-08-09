package types_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/core/types"
	"math/big"
	"testing"
)

func TestEncodeDecodeSignature(t *testing.T) {
	fakeV := new(big.Int).SetUint64(uint64(gofakeit.Uint8()))
	fakeR := big.NewInt(gofakeit.Int64())
	fakeS := big.NewInt(gofakeit.Int64())

	sig := types.NewSignature(fakeV, fakeR, fakeS)
	res, err := types.EncodeSignature(sig)
	Nil(t, err)

	decodedSig, err := types.DecodeSignature(res)
	Nil(t, err)

	Equal(t, sig.V(), decodedSig.V())
	Equal(t, sig.R(), decodedSig.R())
	Equal(t, sig.S(), decodedSig.S())
}

func TestEncodeDecodeAttestation(t *testing.T) {
	domain := gofakeit.Uint32()
	nonce := gofakeit.Uint32()
	root := common.BigToHash(new(big.Int).SetUint64(gofakeit.Uint64()))

	formattedData, err := types.EncodeAttestation(types.NewAttestation(domain, nonce, root))
	Nil(t, err)

	decodedAttestation, err := types.DecodeAttestation(formattedData)
	Nil(t, err)

	Equal(t, decodedAttestation.Domain(), domain)
	Equal(t, decodedAttestation.Nonce(), nonce)

	rawRoot := decodedAttestation.Root()
	Equal(t, common.BytesToHash(rawRoot[:]), root)
}

func TestEncodeDecodeSignedAttestation(t *testing.T) {
	domain := gofakeit.Uint32()
	nonce := gofakeit.Uint32()
	root := common.BigToHash(new(big.Int).SetUint64(gofakeit.Uint64()))

	fakeV := new(big.Int).SetUint64(uint64(gofakeit.Uint8()))
	fakeR := big.NewInt(gofakeit.Int64())
	fakeS := big.NewInt(gofakeit.Int64())

	signedAttestation := types.NewSignedAttestation(
		types.NewAttestation(domain, nonce, root),
		types.NewSignature(fakeV, fakeR, fakeS),
	)

	encoded, err := types.EncodeSignedAttestation(signedAttestation)
	Nil(t, err)

	decoded, err := types.DecodeSignedAttestation(encoded)
	Nil(t, err)

	Equal(t, decoded.Attestation().Root(), signedAttestation.Attestation().Root())
	Equal(t, decoded.Attestation().Domain(), signedAttestation.Attestation().Domain())
	Equal(t, decoded.Attestation().Nonce(), signedAttestation.Attestation().Nonce())

	Equal(t, decoded.Signature().V(), signedAttestation.Signature().V())
	Equal(t, decoded.Signature().R(), signedAttestation.Signature().R())
	Equal(t, decoded.Signature().S(), signedAttestation.Signature().S())
}

func TestEncodeDecodeTips(t *testing.T) {
	// we want to make sure we can deal w/ overflows
	updaterTip := randomUint96BigInt(t)
	relayerTip := randomUint96BigInt(t)
	proverTip := randomUint96BigInt(t)
	processorTip := randomUint96BigInt(t)

	encodedTips, err := types.EncodeTips(types.NewTips(updaterTip, relayerTip, proverTip, processorTip))
	Nil(t, err)

	decodedTips, err := types.DecodeTips(encodedTips)
	Nil(t, err)

	Equal(t, decodedTips.UpdaterTip(), updaterTip)
	Equal(t, decodedTips.RelayerTip(), relayerTip)
	Equal(t, decodedTips.ProverTip(), proverTip)
	Equal(t, decodedTips.ProcessorTip(), processorTip)
}

func TestNewMessageEncodeDecode(t *testing.T) {
	origin := gofakeit.Uint32()
	sender := common.BigToHash(big.NewInt(gofakeit.Int64()))
	nonce := gofakeit.Uint32()
	destination := gofakeit.Uint32()
	body := []byte(gofakeit.Sentence(gofakeit.Number(5, 15)))
	recipient := common.BigToHash(big.NewInt(gofakeit.Int64()))
	optimisticSeconds := gofakeit.Uint32()

	header := types.NewHeader(
		origin, sender, nonce, destination, recipient, optimisticSeconds)

	updaterTip := randomUint96BigInt(t)
	relayerTip := randomUint96BigInt(t)
	proverTip := randomUint96BigInt(t)
	processorTip := randomUint96BigInt(t)

	tips := types.NewTips(updaterTip, relayerTip, proverTip, processorTip)

	newMessage := types.NewMessage(header, tips, body)

	Equal(t, newMessage.OriginDomain(), origin)
	Equal(t, newMessage.Sender(), sender)
	Equal(t, newMessage.Nonce(), nonce)
	Equal(t, newMessage.DestinationDomain(), destination)
	Equal(t, newMessage.Body(), body)

	encodedMessage, err := types.EncodeMessage(newMessage)
	Nil(t, err)

	// make sure decode is same as encode
	decodedMessage, err := types.DecodeMessage(encodedMessage)
	Nil(t, err)

	Equal(t, newMessage.OriginDomain(), decodedMessage.OriginDomain())
	Equal(t, newMessage.Sender(), decodedMessage.Sender())
	Equal(t, newMessage.Nonce(), decodedMessage.Nonce())
	Equal(t, newMessage.DestinationDomain(), decodedMessage.DestinationDomain())
	Equal(t, newMessage.Body(), decodedMessage.Body())
}

func TestHeaderEncodeDecode(t *testing.T) {
	domain := gofakeit.Uint32()
	sender := common.BigToHash(big.NewInt(gofakeit.Int64()))
	nonce := gofakeit.Uint32()
	destination := gofakeit.Uint32()
	recipient := common.BigToHash(big.NewInt(gofakeit.Int64()))
	optimisticSeconds := gofakeit.Uint32()

	ogHeader, err := types.EncodeHeader(types.NewHeader(domain, sender, nonce, destination, recipient, optimisticSeconds))
	Nil(t, err)

	decodedHeader, err := types.DecodeHeader(ogHeader)
	Nil(t, err)

	Equal(t, decodedHeader.OriginDomain(), domain)
	Equal(t, decodedHeader.Sender(), sender)
	Equal(t, decodedHeader.Nonce(), nonce)
	Equal(t, decodedHeader.DestinationDomain(), destination)
	Equal(t, decodedHeader.Recipient(), recipient)
	Equal(t, decodedHeader.OptimisticSeconds(), optimisticSeconds)
}