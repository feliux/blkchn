package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/feliux/blkchn/blockchain"
	"github.com/feliux/blkchn/transaction"
	"github.com/feliux/blkchn/utils"
	"github.com/feliux/blkchn/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port int
}

func NewBlockchainServer(port int) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() int {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("miner_blockchain_address: %v | private_key: %v | publick_key: %v", minersWallet.BlockchainAddress(), minersWallet.PrivateKeyStr(), minersWallet.PublicKeyStr())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		response, err := bc.MarshalJSON()
		if err != nil {
			log.Printf("ERROR marshaling data: %s" + err.Error())
		}
		io.WriteString(w, string(response[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method.")

	}
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet: // get transactions
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		response, err := json.Marshal(struct {
			Transactions []*transaction.Transaction `json:"transactions"`
			Length       int                        `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		if err != nil {
			log.Printf("ERROR marshaling data: %s" + err.Error())
		}
		io.WriteString(w, string(response[:]))

	case http.MethodPost: // create transaction
		decoder := json.NewDecoder(req.Body)
		var t transaction.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR decoding data: %s" + err.Error())
			io.WriteString(w, string(utils.JsonStatus("failed")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("failed")))
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bcs.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)
		w.Header().Add("Content-Type", "application/json")
		var response []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			response = utils.JsonStatus("failed")
		} else {
			w.WriteHeader(http.StatusCreated)
			response = utils.JsonStatus("success")
		}
		io.WriteString(w, string(response))
	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)
		var t transaction.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR decoding data: %s" + err.Error())
			io.WriteString(w, string(utils.JsonStatus("failed")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("failed")))
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bcs.GetBlockchain()
		// AddTransaction not CreateTransaction. We want to sync between nodes
		isUpdated := bc.AddTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)
		w.Header().Add("Content-Type", "application/json")
		var response []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			response = utils.JsonStatus("failed")
		} else {
			response = utils.JsonStatus("success")
		}
		io.WriteString(w, string(response))
	case http.MethodDelete:
		bc := bcs.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JsonStatus("success")))
	default:
		log.Println("ERROR: Invalid HTTP Method.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bcs.GetBlockchain()
		isMined := bc.Mining()

		var response []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			response = utils.JsonStatus("failed")
		} else {
			response = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(response))
	default:
		log.Println("ERROR: Invalid HTTP Method.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) StartMine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bcs.GetBlockchain()
		bc.StartMining()

		response := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(response))
	default:
		log.Println("ERROR: Invalid HTTP Method.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		amount := bcs.GetBlockchain().CalculateTotalAmount(blockchainAddress)
		amountResponse := &blockchain.AmountResponse{Amount: amount}
		response, err := amountResponse.MarshalJSON()
		if err != nil {
			log.Printf("ERROR marshaling data: %s" + err.Error())
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(response[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Consensus(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		bc := bcs.GetBlockchain()
		replaced := bc.ResolveConflicts()
		w.Header().Add("Content-Type", "application/json")
		if replaced {
			io.WriteString(w, string(utils.JsonStatus("success")))
		} else {
			io.WriteString(w, string(utils.JsonStatus("failed")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	bcs.GetBlockchain().Run()
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/mine", bcs.Mine)
	// we can simulate multiple nodes trying to mine with multiples call to localhost:5000/mine/start
	http.HandleFunc("/mine/start", bcs.StartMine)
	http.HandleFunc("/amount", bcs.Amount) // localhost:5000/amount?blockchain_address=XXXXXX
	http.HandleFunc("/consensus", bcs.Consensus)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(bcs.Port()), nil))
}
