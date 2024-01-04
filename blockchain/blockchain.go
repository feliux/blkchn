package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/feliux/blkchn/signature"
	"github.com/feliux/blkchn/transaction"
	"github.com/feliux/blkchn/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 15

	NEIGHBOR_IP_RANGE_START           = 0    // not included
	NEIGHBOR_IP_RANGE_END             = 1    // included
	BLOCKCHAIN_PORT_RANGE_START       = 5000 // included
	BLOCKCHAIN_PORT_RANGE_END         = 5003 // included
	BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC = 20
)

func NewBlockchain(blockchainAddress string, port int) *Blockchain {
	b := Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) Chain() []*Block {
	// to review: to make it public (on ResolveConflicts)
	return bc.chain
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
	// ResolveConflicts when a node is addded
	bc.ResolveConflicts() // to review: if a new node has a malicious valid longestChain? 51% attack
	bc.StartMining()
}

func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(),
		bc.port,
		NEIGHBOR_IP_RANGE_START,
		NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START,
		BLOCKCHAIN_PORT_RANGE_END)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *Blockchain) TransactionPool() []*transaction.Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	str := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*transaction.Transaction{}
	// when a block is created we must delete transactions from other nodes
	for _, hostAddr := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", hostAddr)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		_, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR sending DELETE to blockchain server: %s" + err.Error())
		}
		log.Printf("Deleted transactions on node %s", hostAddr)
	}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for n, block := range bc.chain {
		block.Print(n)
	}
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *signature.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	if isTransacted {
		// replicate transaction information on neighbors
		for _, hostAddr := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &transaction.TransactionRequest{&sender, &recipient, &publicKeyStr, &value, &signatureStr}
			m, err := json.Marshal(bt)
			if err != nil {
				log.Printf("ERROR marshaling data: %s" + err.Error())
			}
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", hostAddr)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			_, err = client.Do(req)
			if err != nil {
				log.Printf("ERROR sending PUT transactions to blockchain server: %s" + err.Error())
				return false
			}
			log.Printf("Updated transaction on node %s", hostAddr)
		}
	}

	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *signature.Signature) bool {
	t := transaction.NewTransaction(sender, recipient, value)
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		if bc.CalculateTotalAmount(sender) < value {
			log.Printf("ERROR: Not enough balance in the wallet %s", sender)
			return false
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Printf("ERROR: could not verify transaction for wallet %s", sender)
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *signature.Signature, t *transaction.Transaction) bool {
	m, err := json.Marshal(t)
	if err != nil {
		log.Printf("ERROR marshaling data: %s" + err.Error())
	}
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*transaction.Transaction {
	transactions := make([]*transaction.Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(
			transactions, transaction.NewTransaction(
				t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.Value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*transaction.Transaction, difficulty int) bool {
	// rule: if new block hash starts by 'difficulty' times 0
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		timestamp:    0, // does not matter. Only matters the rule
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock() // the first routine perform mining
	defer bc.mux.Unlock()
	// in btc it is possible to mine and add a empty block after 10min
	// the probability is close to zero
	/*if len(bc.transactionPool) == 0 {
		log.Println("No transactions to add. Not mining...")
		return false
	}*/

	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("Mining block for current transactions...")
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		_, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR sending PUT consensus to blockchain server: %s" + err.Error())
		}
	}
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	// waits for the duration to elapse and then calls f in its own goroutine (multiples nodes)
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.Value
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	preBlock := chain[0] // genesis block
	currentIndex := 1
	// check if previous hash block matches
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			log.Println("Invalid chain: previousHash is not equal to calculated previous hash.")
			return false
		}
		if !bc.ValidProof(b.Nonce(), b.PreviousHash(), b.Transactions(), MINING_DIFFICULTY) {
			log.Println("Invalid chain: can not proobe difficulty.")
			return false
		}
		preBlock = b
		currentIndex += 1
	}
	return true
}

func (bc *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.chain)
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Printf("ERROR retrieving chain from host %s: %s", endpoint, err.Error())
			// to review: what happen if the node is not reacheable?
			//return false
		}
		if resp.StatusCode == 200 {
			var bcResponse Blockchain
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&bcResponse)
			if err != nil {
				log.Printf("ERROR decoding data verifying longestChain: %s" + err.Error())
				// to review: what happen if the node is not reacheable?
				//return false
			}
			chainResponse := bcResponse.Chain()
			if len(chainResponse) > maxLength && bc.ValidChain(chainResponse) {
				maxLength = len(chainResponse)
				longestChain = chainResponse
			}
		}
	}
	if longestChain != nil {
		bc.chain = longestChain
		log.Println("The chain has been replaced by the longestChain.")
		return true
	}
	// to review
	log.Println("The chain NOT has been replaced by the longestChain. The current chain is the valid chain.")
	return false
}
