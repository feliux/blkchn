package walletserver

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/feliux/blkchn/transaction"
	"github.com/feliux/blkchn/utils"
	"github.com/feliux/blkchn/wallet"
)

const tempDir = "../walletserver/templates"

type WalletServer struct {
	port    int
	gateway string
}

func NewWalletServer(port int, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() int {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Println("ERROR: Invalid HTTP Method.")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, err := myWallet.MarshalJSON()
		if err != nil {
			log.Printf("ERROR marshaling data: %s" + err.Error())
		}
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method.")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// decode the request body into json
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
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
		// convert string data format that can be handled by go
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey) // 128 bytes
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("ERROR: parsing TransactionRequest value.")
			io.WriteString(w, string(utils.JsonStatus("failed")))
			return
		}
		value32 := float32(value)
		walletTransaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := walletTransaction.GenerateSignature()
		signatureStr := signature.String()
		// pending to move to blockchain package ?
		// do not use bt = t as TransactionRequest cause t is the body transaction without the signature
		// to review...
		bt := &transaction.TransactionRequest{
			t.SenderBlockchainAddress,
			t.RecipientBlockchainAddress,
			t.SenderPublicKey,
			&value32,
			&signatureStr,
		}
		m, err := json.Marshal(bt)
		if err != nil {
			log.Printf("ERROR marshaling data: %s" + err.Error())
		}
		buf := bytes.NewBuffer(m)
		w.Header().Add("Content-Type", "application/json")
		resp, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if err != nil {
			log.Printf("ERROR sending transaction to blockchain nodes: %s" + err.Error())
		}
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("failed")))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method.")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(ws.Port()), nil))
}
