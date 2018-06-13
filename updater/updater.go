package updater

import (
	"net/http"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"crypto/ecdsa"
	"encoding/json"
	"encoding/hex"
	"math/big"
	"errors"
)

const (
	updaterURL = "https://c.netlify.com/latest.version"
	signatureURL = "https://c.netlify.com/signature.json"
	currentVersion = "alphax1"
	publicKey = "-----BEGIN PUBLIC KEY-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAAL3kxinRmcZ/mfGZXJakT/J+GwMF zRUW6IA36BiT10xgTt9nhK2GvXADL9goAqO5c7UnoQhb08d61+K2sH7WHkUBCmUJ\nk7v83YRymbemymHdXcMsoVJZ8UxXP1cduuxxCONlO2GDKg5lyB/sDZ56hWkhXIah\nm1NaajeU3j+mHOuo0E4=\n-----END PUBLIC KEY-----"
	securityBreachError = "failed to parse verification data, this is likely a security breach. email contact@larry.science about this"
)



type verificationData struct {
	R string
	S string
	Sum string
}
func getHTTP (url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil  {
		return "", err
	}
	return string(body), nil
}
func Check() (bool, error) {
	latestVersion, err := getHTTP(updaterURL)
	if err != nil {
		return false, errors.New("failed to check for updates, do you have internet?" + err.Error())
	}
	if latestVersion == currentVersion {
		return false, nil
	}
	signatureData, err := getHTTP(signatureURL)
	if err != nil {
		return false, errors.New("failed to get verification data, this is likely a security breach. email contact@larry.science " + err.Error())
	}
	var signature verificationData
	err = json.Unmarshal([]byte(signatureData), &signature)
	if err != nil {
		return false, errors.New(securityBreachError + err.Error())
	}
	block, _ := pem.Decode([]byte(publicKey))
	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	sum, err := hex.DecodeString(signature.Sum)
	if err != nil {
		return false, errors.New(securityBreachError + err.Error())
	}
	var r big.Int
	r.SetString(signature.R, 10)
	var s big.Int
	s.SetString(signature.S, 10)
	if ecdsa.Verify(publicKey.(*ecdsa.PublicKey), sum, &r, &s) == false {
		return false, errors.New(securityBreachError)

	}
	return true, nil
}