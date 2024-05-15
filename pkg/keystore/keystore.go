package keystore

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"os"
)

var pswCache = make(map[string][]byte)

func KeypairFromEth(path string) (*keystore.Key, error) {
	// Make sure key exists before prompting password
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file not found: %s", path)
	}

	var pswd = pswCache[path]
	if len(pswd) == 0 {
		pswd = GetPassword(fmt.Sprintf("Enter password for key %s:", path))
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keyFile failed, err:%s", err)
	}
	ret, err := keystore.DecryptKey(file, string(pswd))
	if err != nil {
		return nil, fmt.Errorf("DecryptKey failed, err:%s", err)
	}
	pswCache[path] = pswd

	return ret, nil
}
