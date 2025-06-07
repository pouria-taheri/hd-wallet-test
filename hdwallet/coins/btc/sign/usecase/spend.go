package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/btcd/txauthor"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/btcd/txscript"
	bd "git.mazdax.tech/blockchain/hdwallet/coins/btc/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// SpendWitnessKeyHash generates, and sets a valid witness for spending the
// passed pkScript with the specified input amount. The input amount *must*
// correspond to the output value of the previous pkScript, or else verification
// will fail since the new sighash digest algorithm defined in BIP0143 includes
// the input value in the sighash.
func (uc *signer) SpendWitnessKeyHash(
	inputAddress *bd.Address, chainParams *chaincfg.Params,
	hashCache *txscript.TxSigHashes, amt int64, tx *wire.MsgTx, idx int,
	secrets txauthor.SecretsSource) error {

	privKey, compressed, err := secrets.GetKey(inputAddress)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendWitnessKeyHash",
			"error":   err,
		}).ErrorF("cannot get private key from input address")
		return err
	}
	pubKey := privKey.PubKey()

	// Once we have the key pair, generate a p2wkh address type, respecting
	// the compression type of the generated key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		uc.logger.With(logger.Field{
			"section":     "SpendWitnessKeyHash",
			"chain param": chainParams,
			"error":       err,
		}).ErrorF("cannot get AddressWitnessPubKeyHash from public key")
		return err
	}

	// With the concrete address type, we can now generate the
	// corresponding witness program to be used to generate a valid witness
	// which will allow us to spend this output.
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendWitnessKeyHash",
			"address": p2wkhAddr.String(),
			"error":   err,
		}).ErrorF("cannot get pkScript from btcutil address")
		return err
	}
	witnessScript, err := txscript.WitnessSignature(hashCache, amt, tx, idx,
		witnessProgram, bd.SigHashAll,
		privKey, true)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendWitnessKeyHash",
			"error":   err,
		}).ErrorF("error in witness signature of unsigned transaction")
		return err
	}

	tx.TxIn[idx].Witness = witnessScript

	return nil
}

// SpendNestedWitnessPubKeyHash generates both a sigScript, and valid witness for
// spending the passed pkScript with the specified input amount. The generated
// sigScript is the version 0 p2wkh witness program corresponding to the queried
// key. The witness stack is identical to that of one which spends a regular
// p2wkh output. The input amount *must* correspond to the output value of the
// previous pkScript, or else verification will fail since the new sighash
// digest algorithm defined in BIP0143 includes the input value in the sighash.
func (uc *signer) SpendNestedWitnessPubKeyHash(inputAddress *bd.Address,
	chainParams *chaincfg.Params, hashCache *txscript.TxSigHashes, amt int64,
	tx *wire.MsgTx, idx int, secrets txauthor.SecretsSource) error {

	privKey, compressed, err := secrets.GetKey(inputAddress)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendNestedWitnessPubKeyHash",
			"error":   err,
		}).ErrorF("cannot get private key from input address")
		return err
	}
	pubKey := privKey.PubKey()

	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	// Next, we'll generate a valid sigScript that'll allow us to spend
	// the p2sh output. The sigScript will contain only a single push of
	// the p2wkh witness program corresponding to the matching public key
	// of this address.
	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, chainParams)
	if err != nil {
		uc.logger.With(logger.Field{
			"section":     "SpendNestedWitnessPubKeyHash",
			"chain param": chainParams,
			"error":       err,
		}).ErrorF("cannot get AddressWitnessPubKeyHash from public key")
		return err
	}
	witnessProgram, err := txscript.PayToAddrScript(p2wkhAddr)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendNestedWitnessPubKeyHash",
			"address": p2wkhAddr.String(),
			"error":   err,
		}).ErrorF("cannot get pkScript from btcutil address")
		return err
	}
	bldr := txscript.NewScriptBuilder()
	bldr.AddData(witnessProgram)
	sigScript, err := bldr.Script()
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendNestedWitnessPubKeyHash",
			"error":   err,
		}).ErrorF("error in script builder")
		return err
	}

	// With the sigScript in place, we'll next generate the proper witness
	// that'll allow us to spend the p2wkh output.
	witnessScript, err := txscript.WitnessSignature(hashCache, amt, tx, idx,
		sigScript, bd.SigHashAll, privKey, compressed)
	if err != nil {
		uc.logger.With(logger.Field{
			"section": "SpendNestedWitnessPubKeyHash",
			"error":   err,
		}).ErrorF("error in witness signature of unsigned transaction")
		return err
	}

	tx.TxIn[idx].Witness = witnessScript

	return nil
}
