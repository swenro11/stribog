package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"regexp"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"

	"github.com/JulianToledano/goingecko"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// BlockchainService 4 ETH-like blockchains.
type BlockchainService struct {
	tokenInfoRepo TokenInfoRepo
	cgCoinRepo    CgCoinRepo
	log           *log.Logger
	cfg           *config.Config
}

// NewBlockchainService -.
func NewBlockchainService(tir TokenInfoRepo, cgr CgCoinRepo, l *log.Logger, cfg *config.Config) *BlockchainService {
	return &BlockchainService{
		tokenInfoRepo: tir,
		cgCoinRepo:    cgr,
		log:           l,
		cfg:           cfg,
	}
}

// 4 TrustWallet coins
const (
	_Arbitrum  = "ARBITRUM"
	_Avalanche = "AVALANCHE"
	_Aurora    = "AURORA"
	_BSC       = "BEP20"
	_Optimism  = "OPTIMISM"
	_Polygon   = "POLYGON"
)

func (service *BlockchainService) GetBlockchainClient(ctx context.Context, blockchainUrl string) (*ethclient.Client, error) {
	url := blockchainUrl
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// based on https://blog.logrocket.com/ethereum-development-using-go-ethereum/ & https://goethereumbook.org/account-balance/
func (service *BlockchainService) GetEthLikeWalletBlance(ctx context.Context, blockchainUrl string, walletAddress string) (*big.Float, error) {
	url := blockchainUrl
	client, err := service.GetBlockchainClient(ctx, url)
	if err != nil {
		return nil, err
	}

	block, errGetLatestBlock := service.GetEthLikeLatestBlock(ctx, client)
	if errGetLatestBlock != nil {
		return nil, errGetLatestBlock
	}
	service.log.Info("Check blockchain Network, last block number: %w", block.Number())

	address := common.HexToAddress(walletAddress)
	balance, errBalance := client.BalanceAt(ctx, address, nil)
	if errBalance != nil {
		return nil, errBalance
	}

	//return result in *big.Int. Exmple - 7115893362963885266, real balance - 7.115893362963885266 MATIC in Polygon
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	floatValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return floatValue, nil
}

// verify your connection to the ETH-like block node by querying for the current block number of the Polygon blockchain
func (service *BlockchainService) GetEthLikeLatestBlock(ctx context.Context, client *ethclient.Client) (*types.Block, error) {
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// Get all blockchain tokens, based on https://github.com/niccoloCastelli/defiarb/blob/master/cmd/updateTokens.go
func (service *BlockchainService) UpdateEthLikeTokens(ctx context.Context) (int, error) {
	fmt.Println("Cloning git repository...")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/trustwallet/assets.git",
	})
	if err != nil {
		return 0, err
	}
	head, err := r.Head()
	if err != nil {
		return 0, err
	}
	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		return 0, err
	}
	fileIter, err := commit.Files()
	if err != nil {
		return 0, err
	}
	fmt.Println("Searching files...")
	//for BSC
	matchAssetReBSC, err := regexp.Compile(`^blockchains/smartchain/assets/([0-9A-Za-z]+)/info\.json$`)
	if err != nil {
		return 0, err
	}
	//for Polygon
	matchAssetRePolygon, err := regexp.Compile(`^blockchains/polygon/assets/([0-9A-Za-z]+)/info\.json$`)
	if err != nil {
		return 0, err
	}
	//for Arbitrum
	matchAssetReArbitrum, err := regexp.Compile(`^blockchains/arbitrum/assets/([0-9A-Za-z]+)/info\.json$`)
	if err != nil {
		return 0, err
	}
	//for AvalancheC
	matchAssetReAvalanche, err := regexp.Compile(`^blockchains/avalanchec/assets/([0-9A-Za-z]+)/info\.json$`)
	if err != nil {
		return 0, err
	}
	//for Optimism
	matchAssetReOptimism, err := regexp.Compile(`^blockchains/optimism/assets/([0-9A-Za-z]+)/info\.json$`)
	if err != nil {
		return 0, err
	}
	err = fileIter.ForEach(func(f *object.File) error {
		savePolygonErr := service.SaveToken(ctx, matchAssetRePolygon, f, _Polygon)
		if savePolygonErr != nil {
			return savePolygonErr
		}
		saveBSCerr := service.SaveToken(ctx, matchAssetReBSC, f, _BSC)
		if saveBSCerr != nil {
			return saveBSCerr
		}
		saveArbitrumErr := service.SaveToken(ctx, matchAssetReArbitrum, f, _Arbitrum)
		if saveArbitrumErr != nil {
			return saveArbitrumErr
		}
		saveAvalancheErr := service.SaveToken(ctx, matchAssetReAvalanche, f, _Avalanche)
		if saveAvalancheErr != nil {
			return saveAvalancheErr
		}
		saveOptimismErr := service.SaveToken(ctx, matchAssetReOptimism, f, _Optimism)
		if saveOptimismErr != nil {
			return saveOptimismErr
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	allTokens, err := service.tokenInfoRepo.GetAllTokens(ctx)
	if err != nil {
		return 0, err
	}
	return len(allTokens), nil
}

func (service *BlockchainService) SaveToken(ctx context.Context, blockchainRegexp *regexp.Regexp, f *object.File, blockchain string) error {
	m := blockchainRegexp.FindStringSubmatch(f.Name)
	if len(m) == 0 {
		return nil
	}
	content, err := f.Contents()
	if err != nil {
		return err
	}

	tokenInfo := entity.TokenInfo{}
	if err := json.Unmarshal([]byte(content), &tokenInfo); err != nil {
		return err
	}
	if tokenInfo.Type == blockchain && tokenInfo.Status == "active" {
		existTokens, errExistToken := service.tokenInfoRepo.GetByIdAndType(ctx, blockchain, tokenInfo.Id)
		if errExistToken != nil {
			return errExistToken
		}
		if len(existTokens) > 0 {
			errUti := service.tokenInfoRepo.UpdateTokenInfo(ctx, existTokens[0], tokenInfo)
			if errUti != nil {
				return errUti
			}
		} else {
			errSti := service.tokenInfoRepo.StoreTokenInfo(ctx, tokenInfo)
			if errSti != nil {
				return errSti
			}
		}

	}
	return nil
}

func (service *BlockchainService) GetTokenUsdPrice(cgId string) (float64, error) {
	cgClient := goingecko.NewClient(nil)
	defer cgClient.Close()

	data, err := cgClient.CoinsId(cgId, true, true, true, false, false, false)
	if err != nil {
		return 0, err
	}

	return data.MarketData.CurrentPrice.Usd, nil
}

func (service *BlockchainService) CgCoinsList(ctx context.Context) error {
	cgClient := goingecko.NewClient(nil)
	defer cgClient.Close()

	coinsList, errCoinsList := cgClient.CoinsList()
	if errCoinsList != nil {
		return errCoinsList
	}

	errDeleteAll := service.cgCoinRepo.DeleteAll(ctx)
	if errDeleteAll != nil {
		return fmt.Errorf("BlockchainService - CgCoinsList - cgCoinRepo.DeleteAll: ", errDeleteAll)
	}

	for _, cgCoin := range coinsList {
		cgc := entity.CgCoin{}
		cgc.ID = cgCoin.ID
		cgc.Symbol = cgCoin.Symbol
		cgc.Name = cgCoin.Name

		errStore := service.cgCoinRepo.StoreCgCoin(ctx, cgc)
		if errStore != nil {
			return fmt.Errorf("BlockchainService - CgCoinsList - cgCoinRepo.StoreCgCoin: %w", errStore)
		}
	}

	return nil
}

// based on https://github.com/huahuayu/go-transaction-decoder/blob/master/main.go
func (service *BlockchainService) DecodeTxInput(abiReader io.Reader, txInput string) (map[string]interface{}, error) {
	// load contract ABI
	abi, err := abi.JSON(abiReader)
	if err != nil {
		return nil, err
	}

	// decode txInput method signature
	decodedSig, err := hex.DecodeString(txInput[2:10])
	if err != nil {
		return nil, err
	}

	// recover Method from signature and ABI
	method, err := abi.MethodById(decodedSig)
	if err != nil {
		return nil, err
	}

	// decode txInput Payload
	decodedData, err := hex.DecodeString(txInput[10:])
	if err != nil {
		return nil, err
	}

	// unpack method inputs
	inputMap := make(map[string]interface{}, 0)
	err = method.Inputs.UnpackIntoMap(inputMap, decodedData)
	if err != nil {
		return nil, err
	}

	return inputMap, nil
}
