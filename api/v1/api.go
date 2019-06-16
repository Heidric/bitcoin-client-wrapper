package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

type transactionDetails struct {
	Account   string
	Address   string
	Category  string
	Amount    float64
	Label     string
	Vout      int
	Fee       float64
	Abandoned bool
}

type transactionData struct {
	Amount            float64
	Fee               float64
	Confirmations     int
	Blockhash         string
	Blockindex        int
	Txid              string
	Time              int64
	Timereceived      int64
	Bip125Replaceable string
	Details           []transactionDetails
	Hex               string
}

type errorData struct {
	Code    int
	Message string
}

type rpcResponse struct {
	Result transactionData
	Error  errorData
	Id     string
}

type getTransactionArgs struct {
	Txid              string
	Include_watchonly bool
}

type sendToAddressRequestArgs struct {
	Passphrase string
	Timeout    int

	Address               string
	Amount                string
	Comment               string
	Comment_to            string
	Subtractfeefromamount bool
	Replaceable           bool
	Conf_target           int
	Estimate_mode         string
}

type passPhraseArgsStruct struct {
	Passphrase string
	Timeout    int
}

type sendToAddressArgsStruct struct {
	Address               string
	Amount                string
	Comment               string
	Comment_to            string
	Subtractfeefromamount bool
	Replaceable           bool
	Conf_target           int
	Estimate_mode         string
}

type getNewAddressRequestArgs struct {
	Passphrase string
	Timeout    int

	Label        string
	Address_type string
}

type getNewAddressArgsStruct struct {
	Label        string
	Address_type string
}

func jsonResponse(w http.ResponseWriter, data interface{}, code int) {
	if code == 500 && os.Getenv("ENV") == "dev" {
		data = struct {
			Message string
		}{
			Message: "Something went wrong",
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Print(err)
		if code == 500 && os.Getenv("ENV") == "dev" {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonData)

	if err != nil {
		log.Print(err)
		if code == 500 && os.Getenv("ENV") == "dev" {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
		return
	}
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactionId := chi.URLParam(r, "id")

	watchOnly, err := strconv.ParseBool(r.URL.Query().Get("watchonly"))

	if err != nil {
		watchOnly = false
	}

	client, err := jsonrpc.Dial("tcp", os.Getenv("RPC_ADDR"))

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error in dialing. %s", err)
		return
	}

	args := &getTransactionArgs{
		Txid:              transactionId,
		Include_watchonly: watchOnly,
	}

	var result rpcResponse

	err = client.Call("gettransaction", args, &result)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error getting transaction: %d", err)
		return
	}

	err = client.Close()

	if err != nil {
		log.Printf("Error closing connection: %d", err)
	}

	jsonResponse(w, result.Result, http.StatusOK)
}

func SendToAddress(w http.ResponseWriter, r *http.Request) {
	var args sendToAddressRequestArgs

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&args)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error parsing sendToAddress body: %d", err)
		return
	}

	if args.Passphrase == "" {
		requestError := struct {
			Message string
		}{
			Message: "Passphrase is required",
		}
		jsonResponse(w, requestError, http.StatusBadRequest)
		log.Println("Passphrase is required")
		return
	}

	if args.Timeout == 0 {
		args.Timeout = 1
	}

	if args.Estimate_mode == "" {
		args.Estimate_mode = "UNSET"
	}

	if args.Estimate_mode != "UNSET" && args.Estimate_mode != "ECONOMICAL" && args.Estimate_mode != "CONSERVATIVE" {
		requestError := struct {
			Message string
		}{
			Message: "Estimate_mode must be one of the following: UNSET, ECONOMICAL, CONSERVATIVE",
		}
		jsonResponse(w, requestError, http.StatusBadRequest)
		log.Printf("Estimate mode had %s value", args.Estimate_mode)
		return
	}

	client, err := jsonrpc.Dial("tcp", os.Getenv("RPC_ADDR"))

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error in dialing. %s", err)
		return
	}

	var result rpcResponse

	passPhraseArgs := &passPhraseArgsStruct{
		Passphrase: args.Passphrase,
		Timeout:    args.Timeout,
	}

	err = client.Call("walletpassphrase", passPhraseArgs, &result)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error sending Passphrase: %d", err)
		return
	}

	sendToAddressArgs := &sendToAddressArgsStruct{
		Address:               args.Address,
		Amount:                args.Amount,
		Comment:               args.Comment,
		Comment_to:            args.Comment_to,
		Subtractfeefromamount: args.Subtractfeefromamount,
		Replaceable:           args.Replaceable,
		Conf_target:           args.Conf_target,
		Estimate_mode:         args.Estimate_mode,
	}

	err = client.Call("sendtoaddress", sendToAddressArgs, &result)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error getting transaction: %d", err)
		return
	}

	err = client.Close()

	if err != nil {
		log.Printf("Error closing connection: %d", err)
	}

	jsonResponse(w, result.Result, http.StatusOK)
}

func GetNewAddress(w http.ResponseWriter, r *http.Request) {
	var args getNewAddressRequestArgs

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&args)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error parsing sendToAddress body: %d", err)
		return
	}

	if args.Passphrase == "" {
		requestError := struct {
			Message string
		}{
			Message: "Passphrase is required",
		}
		jsonResponse(w, requestError, http.StatusBadRequest)
		log.Println("Passphrase is required")
		return
	}

	if args.Timeout == 0 {
		args.Timeout = 1
	}

	if args.Address_type != "legacy" && args.Address_type != "p2sh-segwit" && args.Address_type != "bech32" {
		requestError := struct {
			Message string
		}{
			Message: "Address_type must be one of the following: legacy, p2sh-segwit, bech32",
		}
		jsonResponse(w, requestError, http.StatusBadRequest)
		log.Printf("Address type had %s value", args.Address_type)
		return
	}

	client, err := jsonrpc.Dial("tcp", os.Getenv("RPC_ADDR"))

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error in dialing. %s", err)
		return
	}

	var result rpcResponse

	passPhraseArgs := &passPhraseArgsStruct{
		Passphrase: args.Passphrase,
		Timeout:    args.Timeout,
	}

	err = client.Call("walletpassphrase", passPhraseArgs, &result)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error getting transaction: %d", err)
		return
	}

	getNewAddressArgs := &getNewAddressArgsStruct{
		Label:        args.Label,
		Address_type: args.Address_type,
	}

	err = client.Call("getnewaddress", getNewAddressArgs, &result)

	if err != nil {
		jsonResponse(w, err, http.StatusInternalServerError)
		log.Printf("Error getting transaction: %d", err)
		return
	}

	err = client.Close()

	if err != nil {
		log.Printf("Error closing connection: %d", err)
	}

	jsonResponse(w, result.Result, http.StatusOK)
}

func BtcRouter() http.Handler {
	router := chi.NewRouter()
	router.Get("/transaction/{id}", GetTransaction)
	router.Post("/send-to-address", SendToAddress)
	router.Post("/get-new-Address", GetNewAddress)
	return router
}
