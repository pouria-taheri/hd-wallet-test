package types

type (
	Block struct {
		Number         string      `json:"number"`
		Hash           string      `json:"hash"`
		ParentHash     string      `json:"parent_hash"`
		StateRoot      string      `json:"state_root"`
		ExtrinsicsRoot string      `json:"extrinsics_root"`
		AuthorId       string      `json:"author_id"`
		Logs           []Log       `json:"logs"`
		Extrinsics     []Extrinsic `json:"extrinsics"`
		Finalized      bool        `json:"finalized"`
	}

	Log struct {
		Type  string   `json:"type"`
		Index string   `json:"index"`
		Value []string `json:"value"`
	}

	Extrinsic struct {
		Method    Method                 `json:"method"`
		Signature Signature              `json:"signature"`
		Nonce     string                 `json:"nonce"`
		Args      map[string]interface{} `json:"args"`
		Hash      string                 `json:"hash"`
		Events    []Event                `json:"events"`
		Success   bool                   `json:"success"`
		PaysFee   bool                   `json:"pays_fee"`
		Info      EstimateFeeResult      `json:"info"`
	}

	Signature struct {
		Signature string `json:"signature"`
		Signer    Signer `json:"pkg"`
	}

	Signer struct {
		ID string `json:"id"`
	}

	Method struct {
		Pallet string `json:"pallet"`
		Method string `json:"method"`
	}

	Event struct {
		Method Method      `json:"method"`
		Data   interface{} `json:"data"`
	}

	ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Stack   string `json:"stack"`
	}

	Head struct {
		Hash          string `json:"hash"`
		Number        string `json:"number"`
		ParentHash    string `json:"parentHash"`
		StateRoot     string `json:"stateRoot"`
		ExtrinsicRoot string `json:"ExtrinsicRoot"`
		Finalized     bool   `json:"finalized"`
	}

	Balance struct {
		TokenSymbol string `json:"tokenSymbol"`
		Free        string `json:"free"`
	}

	BalanceTransferResult struct {
		SigningPayload string `json:"signingPayload"`
		SpecName       string `json:"specName"`
		SpecVersion    int    `json:"specVersion"`
		Method         string `json:"method"`
		Version        int    `json:"version"`
		BlockHash      string `json:"blockHash"`
		Era            string `json:"era"`
		GenesisHash    string `json:"genesisHash"`
		Nonce          string `json:"nonce"`
		Tip            string `json:"tip"`
		MetadataRpc    string `json:"metadataRpc"`
	}

	SubmittableExtrinsic struct {
		Extrinsic string `json:"extrinsic"`
	}

	DryRunResult struct {
		ResultType        string `json:"resultType"`
		Result            string `json:"dryRunResult"`
		ValidityErrorType string `json:"validityErrorType"`
	}

	SubmitExtrinsicResult struct {
		Hash string `json:"hash"`
	}

	EstimateFeeResult struct {
		Weight     string `json:"weight"`
		Class      string `json:"class"`
		PartialFee string `json:"partialFee"`
	}
)
