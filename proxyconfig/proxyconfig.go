package proxyconfig

type ProxyConfig struct {
	Addr           string
	ValidateCert   bool
	EncryptionType string
	BypassType     string
}