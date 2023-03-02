package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wttech/aemc/pkg/common/filex"
	"github.com/wttech/aemc/pkg/common/pathx"
)

type Crypto struct {
	instance *Instance

	keyBundleSymbolicName string
}

func NewCrypto(instance *Instance) *Crypto {
	return &Crypto{
		instance: instance,

		keyBundleSymbolicName: "com.adobe.granite.crypto.file",
	}
}

func (c Crypto) Setup(hmacFile string, masterFile string) (bool, error) {
	if !c.instance.IsLocal() {
		return false, fmt.Errorf("%s > Crypto keys could only be set on local instance", c.instance.ID())
	}
	if !pathx.Exists(hmacFile) {
		return false, fmt.Errorf("%s > Crypto hmac file '%s' does not exist", c.instance.ID(), hmacFile)
	}
	if !pathx.Exists(masterFile) {
		return false, fmt.Errorf("%s > Crypto master file '%s' does not exist", c.instance.ID(), masterFile)
	}
	log.Debugf("%s > setting up Crypto keys (hmac '%s', master '%s')", c.instance.ID(), hmacFile, masterFile)
	osgi := c.instance.OSGI()
	keyBundle, err := osgi.BundleManager().Find(c.keyBundleSymbolicName)
	if err != nil {
		return false, err
	}
	if keyBundle == nil {
		return false, fmt.Errorf("%s > cannot find Crypto key bundle using symbolic name '%s'", c.instance.ID(), c.keyBundleSymbolicName)
	}
	keyDir := fmt.Sprintf("%s/data", c.instance.Local().BundleDir(keyBundle.ID))
	hmacTargetFile := fmt.Sprintf("%s/hmac", keyDir)
	masterTargetFile := fmt.Sprintf("%s/master", keyDir)

	hmacOk, err := filex.Equals(hmacFile, hmacTargetFile)
	if err != nil {
		return false, err
	}
	masterOk, err := filex.Equals(masterFile, masterTargetFile)
	if err != nil {
		return false, err
	}
	if hmacOk && masterOk {
		log.Debugf("%s > skipping setting Crypto keys (hmac '%s', master '%s') as they are up-to-date", c.instance.ID(), hmacFile, masterFile)
		return false, nil
	}
	log.Debugf("%s > copying Crypto hmac file from '%s' to '%s'", c.instance.ID(), hmacFile, hmacTargetFile)
	if err := filex.Copy(hmacFile, hmacTargetFile, true); err != nil {
		return false, fmt.Errorf("%s > cannot copy Crypto hmac file from '%s' to '%s': %w", c.instance.ID(), hmacFile, hmacTargetFile, err)
	}
	log.Debugf("%s > copying Crypto master file from '%s' to '%s'", c.instance.ID(), masterFile, masterTargetFile)
	if err := filex.Copy(masterFile, masterTargetFile, true); err != nil {
		return false, fmt.Errorf("%s > cannot copy Crypto master file from '%s' to '%s'> %w", c.instance.ID(), masterFile, masterTargetFile, err)
	}
	if err := osgi.Restart(); err != nil {
		return false, err
	}
	log.Debugf("%s > set up Crypto keys (hmac '%s', master '%s')", c.instance.ID(), hmacFile, masterFile)
	return true, nil
}

func (c Crypto) Protect(value string) (string, error) {
	return "", nil
}
