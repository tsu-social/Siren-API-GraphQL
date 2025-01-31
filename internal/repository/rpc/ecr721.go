// Package rpc provides high level access to the Fantom Opera blockchain
// node through RPC interface.
package rpc

import (
	"artion-api-graphql/internal/repository/rpc/contracts"
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// defaultMintingTestTokenUrl is the URL we used to test NFT minting calls.
const defaultMintingTestTokenUrl = "https://minter.artion.io/default/access/minter/estimation.json"

// defaultMintingTestFee is the default fee we try on minting test (10 FTM).
var defaultMintingTestFee = hexutil.MustDecodeBig("0x8AC7230489E80000")

// Erc721StartingBlockNumber provides the first important block number for the ERC-721 contract.
// We try to get the first Transfer() event on the contract,
// anything before it is irrelevant for this API.
func (o *Opera) Erc721StartingBlockNumber(adr *common.Address) (uint64, error) {
	// instantiate contract
	erc, err := contracts.NewErc721(*adr, o.ftm)
	log.Debugf("getting ERC-721 starting block number for %s", adr.String())
	if err != nil {
		return 0, err
	}

	// iterate over transfers from zero address (e.g. mint calls)
	iter, err := erc.FilterTransfer(nil, []common.Address{{}}, nil, nil)
	if err != nil {
		return 0, err
	}

	var blk uint64
	if iter.Next() {
		blk = iter.Event.Raw.BlockNumber
	}

	if err := iter.Close(); err != nil {
		log.Errorf("could not close filter iterator; %s", err.Error())
	}
	return blk, nil
}

// CanMintErc721 checks if the given user can mint a new token on the given NFT contract.
func (o *Opera) CanMintErc721(contract *common.Address, user *common.Address, fee *big.Int) (bool, error) {
	data, err := o.abiFantom721.Pack("mint", *user, defaultMintingTestTokenUrl)
	if err != nil {
		return false, err
	}

	// use default fee, if not specified
	if fee == nil {
		fee = o.MustPlatformFee(contract)
		log.Infof("platform fee for %s is %s", contract.String(), (*hexutil.Big)(fee).String())
	}

	// try to estimate the call
	gas, err := o.ftm.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  *user,
		To:    contract,
		Data:  data,
		Value: fee,
	})
	if err != nil {
		log.Warningf("user %s can not mint on ERC-721 %s; %s", user.String(), contract.String(), err.Error())
		return false, nil
	}

	log.Infof("user %s can mint on ERC-721 %s for %d gas", user.String(), contract.String(), gas)
	return true, nil
}

// MustPlatformFee returns the platform fee for the given contract, or the default one.
func (o *Opera) MustPlatformFee(contract *common.Address) *big.Int {
	data, err := o.ftm.CallContract(context.Background(), ethereum.CallMsg{
		From: common.Address{},
		To:   contract,
		Data: common.Hex2Bytes("26232a2e"),
	}, nil)
	if err != nil {
		log.Errorf("can not get platform fee from %s; %s", contract.String(), err.Error())
		return defaultMintingTestFee
	}

	// try to unpack the data if possible; we expect uint256 value = 32 bytes
	if len(data) != 32 {
		log.Errorf("invalid platform fee response from %s; expected 32 bytes, %d bytes received", contract.String(), len(data))
		return defaultMintingTestFee
	}

	return new(big.Int).SetBytes(data)
}

// Erc721TokenUri gets a token specific URI address from ERC-721 contract using tokenURI() call.
func (o *Opera) Erc721TokenUri(contract *common.Address, tokenId *big.Int) (string, error) {
	// prepare params
	input, err := o.Erc721Abi().Pack("tokenURI", tokenId)
	if err != nil {
		log.Errorf("can not pack data; %s", err.Error())
		return "", err
	}

	// call the contract
	data, err := o.ftm.CallContract(context.Background(), ethereum.CallMsg{
		From: common.Address{},
		To:   contract,
		Data: input,
	}, nil)
	res, err := o.abiFantom721.Unpack("tokenURI", data)
	if err != nil {
		log.Errorf("can not decode response; %s", err.Error())
		return "", err
	}

	return strings.Replace(
		*abi.ConvertType(res[0], new(string)).(*string),
		"{id}",
		fmt.Sprintf("%064x", tokenId), -1), nil
}

func (o *Opera) Erc721IsApprovedForAll(contract *common.Address, owner *common.Address, operator *common.Address) (bool, error) {
	// instantiate contract
	erc, err := contracts.NewErc721(*contract, o.ftm)
	if err != nil {
		return false, err
	}
	return erc.IsApprovedForAll(nil, *owner, *operator)
}
