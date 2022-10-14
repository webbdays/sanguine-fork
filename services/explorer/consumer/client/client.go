// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package client

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/synapsecns/sanguine/services/explorer/consumer/client/model"
)

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) *Client {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	Logs                   []*model.Log         "json:\"logs\" graphql:\"logs\""
	LogsRange              []*model.Log         "json:\"logsRange\" graphql:\"logsRange\""
	Receipts               []*model.Receipt     "json:\"receipts\" graphql:\"receipts\""
	ReceiptsRange          []*model.Receipt     "json:\"receiptsRange\" graphql:\"receiptsRange\""
	Transactions           []*model.Transaction "json:\"transactions\" graphql:\"transactions\""
	TransactionsRange      []*model.Transaction "json:\"transactionsRange\" graphql:\"transactionsRange\""
	BlockTime              *int                 "json:\"blockTime\" graphql:\"blockTime\""
	LastStoredBlockNumber  *int                 "json:\"lastStoredBlockNumber\" graphql:\"lastStoredBlockNumber\""
	FirstStoredBlockNumber *int                 "json:\"firstStoredBlockNumber\" graphql:\"firstStoredBlockNumber\""
	TxSender               *string              "json:\"txSender\" graphql:\"txSender\""
}
type GetLogsRange struct {
	Response []*struct {
		ContractAddress string   "json:\"contract_address\" graphql:\"contract_address\""
		ChainID         int      "json:\"chain_id\" graphql:\"chain_id\""
		Topics          []string "json:\"topics\" graphql:\"topics\""
		Data            string   "json:\"data\" graphql:\"data\""
		BlockNumber     int      "json:\"block_number\" graphql:\"block_number\""
		TxHash          string   "json:\"tx_hash\" graphql:\"tx_hash\""
		TxIndex         int      "json:\"tx_index\" graphql:\"tx_index\""
		BlockHash       string   "json:\"block_hash\" graphql:\"block_hash\""
		Index           int      "json:\"index\" graphql:\"index\""
		Removed         bool     "json:\"removed\" graphql:\"removed\""
	} "json:\"response\" graphql:\"response\""
}
type GetBlockTime struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetLastStoredBlockNumber struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetFirstStoredBlockNumber struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetTxSender struct {
	Response *string "json:\"response\" graphql:\"response\""
}

const GetLogsRangeDocument = `query GetLogsRange ($chain_id: Int!, $start_block: Int!, $end_block: Int!, $page: Int!) {
	response: logsRange(chain_id: $chain_id, start_block: $start_block, end_block: $end_block, page: $page) {
		contract_address
		chain_id
		topics
		data
		block_number
		tx_hash
		tx_index
		block_hash
		index
		removed
	}
}
`

func (c *Client) GetLogsRange(ctx context.Context, chainID int, startBlock int, endBlock int, page int, httpRequestOptions ...client.HTTPRequestOption) (*GetLogsRange, error) {
	vars := map[string]interface{}{
		"chain_id":    chainID,
		"start_block": startBlock,
		"end_block":   endBlock,
		"page":        page,
	}

	var res GetLogsRange
	if err := c.Client.Post(ctx, "GetLogsRange", GetLogsRangeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetBlockTimeDocument = `query GetBlockTime ($chain_id: Int!, $block_number: Int!) {
	response: blockTime(chain_id: $chain_id, block_number: $block_number)
}
`

func (c *Client) GetBlockTime(ctx context.Context, chainID int, blockNumber int, httpRequestOptions ...client.HTTPRequestOption) (*GetBlockTime, error) {
	vars := map[string]interface{}{
		"chain_id":     chainID,
		"block_number": blockNumber,
	}

	var res GetBlockTime
	if err := c.Client.Post(ctx, "GetBlockTime", GetBlockTimeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLastStoredBlockNumberDocument = `query GetLastStoredBlockNumber ($chain_id: Int!) {
	response: lastStoredBlockNumber(chain_id: $chain_id)
}
`

func (c *Client) GetLastStoredBlockNumber(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetLastStoredBlockNumber, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetLastStoredBlockNumber
	if err := c.Client.Post(ctx, "GetLastStoredBlockNumber", GetLastStoredBlockNumberDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetFirstStoredBlockNumberDocument = `query GetFirstStoredBlockNumber ($chain_id: Int!) {
	response: firstStoredBlockNumber(chain_id: $chain_id)
}
`

func (c *Client) GetFirstStoredBlockNumber(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetFirstStoredBlockNumber, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetFirstStoredBlockNumber
	if err := c.Client.Post(ctx, "GetFirstStoredBlockNumber", GetFirstStoredBlockNumberDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetTxSenderDocument = `query GetTxSender ($chain_id: Int!, $tx_hash: String!) {
	response: txSender(chain_id: $chain_id, tx_hash: $tx_hash)
}
`

func (c *Client) GetTxSender(ctx context.Context, chainID int, txHash string, httpRequestOptions ...client.HTTPRequestOption) (*GetTxSender, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
		"tx_hash":  txHash,
	}

	var res GetTxSender
	if err := c.Client.Post(ctx, "GetTxSender", GetTxSenderDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}