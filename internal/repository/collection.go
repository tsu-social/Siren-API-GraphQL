// Package repository implements persistent data access and processing.
package repository

import (
	"artion-api-graphql/internal/types"
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// addCollectionQueueCapacity is the capacity of the queue being filled with
// new collections to be added into the repository.
const addCollectionQueueCapacity = 10

// NewCollectionQueue provides queue filled with addresses of collection contracts
// to be added into the API.
func (p *Proxy) NewCollectionQueue() chan common.Address {
	return p.newCollectionQueue
}

// AddNewCollection pushes the given address to be validated and added as a new collection.
func (p *Proxy) AddNewCollection(adr *common.Address) {
	p.newCollectionQueue <- *adr
}

// CollectionName provides a name of an Artion ERC721 and/or ERC1155 token.
func (p *Proxy) CollectionName(adr *common.Address) (string, error) {
	return p.rpc.CollectionName(adr)
}

// CollectionSymbol provides a symbol of an Artion ERC721 and/or ERC1155 token.
func (p *Proxy) CollectionSymbol(adr *common.Address) (string, error) {
	return p.rpc.CollectionSymbol(adr)
}

// StoreCollection adds the specified NFT collection contract record into persistent storage.
func (p *Proxy) StoreCollection(nft *types.Collection) error {
	return p.db.StoreCollection(nft)
}

func (p *Proxy) GetCollection(address common.Address) (*types.Collection, error) {
	key := "Col-" + address.String()
	user, err, _ := p.callGroup.Do(key, func() (interface{}, error) {
		return p.db.GetCollection(address)
	})
	return user.(*types.Collection), err
}

func (p *Proxy) ListCollections(cursor types.Cursor, count int, backward bool) (*types.CollectionList, error) {
	return p.db.ListCollections(cursor, count, backward)
}

// NFTContractType analyses the contract on given address and returns the type, if possible.
func (p *Proxy) NFTContractType(adr *common.Address) (string, error) {
	if p.IsErc721Contract(adr) {
		log.Infof("contract %s is ERC-721", adr.String())
		return types.ContractTypeERC721, nil
	}
	if p.IsErc1155Contract(adr) {
		log.Infof("contract %s is ERC-1155", adr.String())
		return types.ContractTypeERC1155, nil
	}
	return "", fmt.Errorf("unknown contract type at %s", adr.String())
}

// CanMint checks if the given user can mint a new token on the given NFT contract.
func (p *Proxy) CanMint(contract *common.Address, user *common.Address, fee *big.Int) (bool, error) {
	// the ERC-721 minter differs from other contracts, we need to check the type first
	if p.IsErc721Contract(contract) {
		return p.rpc.CanMintErc721(contract, user, fee)
	}

	// it's either ERC-1155 or not a valid minter at all
	return p.rpc.CanMintErc1155(contract, user, fee)
}

// CollectionOwner tries to get the owner of the given collection.
func (p *Proxy) CollectionOwner(contract *common.Address) *common.Address {
	return p.rpc.CollectionOwner(contract)
}

// CanRegisterCollection checks if the given collection can be registered.
func (p *Proxy) CanRegisterCollection(contract *common.Address, user *common.Address) error {
	if !p.IsErc721Contract(contract) {
		return fmt.Errorf("the contract is not ERC-721")
	}

	owner := p.rpc.CollectionOwner(contract)
	if owner != nil {
		if !bytes.Equal(user.Bytes(), owner.Bytes()) {
			return fmt.Errorf("the contract is owned by somebody else")
		}
	}

	return nil
}

func (p *Proxy) IsApprovedForAll(contract *common.Address, owner *common.Address, operator *common.Address) (bool, error) {
	isApproved, err := p.rpc.Erc721IsApprovedForAll(contract, owner, operator)
	if err != nil {
		return false, fmt.Errorf("IsApprovedForAll failed; %s", err)
	}
	return isApproved, err
}
