# Bitcoin-client-wrapper

## Launch instructions

1. Set two required global variables: RPC_ADDR (address of the bitcoin node to connect to) and MAIN_PORT (port to run the server on)
2. Set the optional global variable, which is responsible for the debug log in the server responses: ENV. Set it to "dev" if you want the extra debug info in case of errors
3. Either build and run source files or launch the built binary

## API description

All the methods return JSON object

<pre>
{
  "Message": "Something went wrong"
}
</pre>

with status code 500 in case of the server error.

### Get transaction by it's id
GET /api/v1/transaction/{id}?watchonly=false/true

where {id} should be swapped with the transaction id and the "watchonly" can be provided as the query parameter, which defaults to false if it was not provided

Returns JSON object in the form of

<pre>
{
  "amount": double,
  "fee": double,
  "confirmations" : int,
  "blockhash" : string,
  "blockindex" : int,
  "blocktime" : int,
  "txid" : string,
  "time" : int,
  "timereceived" : int,
  "bip125-replaceable": string,
  "details" : [
    {
      "account" : string,
      "address" : string,
      "category" : string,
      "amount" : float,
      "label" : string,
      "vout" : int,
      "fee": float,
      "abandoned": bool
    }
    ,...
  ],
  "hex" : string
}
</pre>

with the status code 200 if the request was successfully processed

### Get new address
POST /api/v1/get-new-address

Accepts JSON object in the request body, in the following form:

<pre>
{
  "Passphrase": string, required,
  "Timeout": int, optional, defaults to 1,
  "Label": string, optional,
  "Address_type": string, optional, one of the “legacy”, “p2sh-segwit”, and “bech32”
}
</pre>

Returns JSON object in the form of

<pre>
{
  "result": string
}
</pre>

with the status code 200 if the request was successfully processed

Returns

<pre>
{
  "Message": "Passphrase is required"
}
</pre>

with status code 400 if the passphrase wasn't provided

Returns

<pre>
{
  "Message": "Address_type must be one of the following: legacy, p2sh-segwit, bech32"
}
</pre>

with status code 400 if the address type had wrong value

### Send btc to address
POST /api/v1/send-to-address

Accepts JSON object in the request body, in the following form:

<pre>
{
  "Passphrase": string, required,
  "Timeout": int, optional, defaults to 1,

  "Address": string, required
  "Amount": string, required
  "Comment": string, optional
  "Comment_to": string, optional 
  "Subtractfeefromamount": bool, optional
  "Replaceable": bool, optional
  "Conf_target": int, optional
  "Estimate_mode": string, optional, one of the following: "UNSET", "ECONOMICAL", "CONSERVATIVE", defaults to "UNSET"
}
</pre>

Returns JSON object in the form of

<pre>
{
  "txid": string
}
</pre>

with the status code 200 if the request was successfully processed

Returns

<pre>
{
  "Message": "Passphrase is required"
}
</pre>

with status code 400 if the passphrase wasn't provided

Returns

<pre>
{
  "Message": "Estimate_mode must be one of the following: UNSET, ECONOMICAL, CONSERVATIVE"
}
</pre>

with status code 400 if the estimate mode had wrong value

