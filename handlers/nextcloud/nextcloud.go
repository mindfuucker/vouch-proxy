package nextcloud

import (
	"encoding/json"
	"github.com/vouch/vouch-proxy/handlers/common"
	"github.com/vouch/vouch-proxy/pkg/cfg"
	"github.com/vouch/vouch-proxy/pkg/structs"
	"io/ioutil"
	"net/http"
)

type Handler struct{}

var (
	log = cfg.Cfg.Logger
)

func (Handler) GetUserInfo(r *http.Request, user *structs.User, customClaims *structs.CustomClaims, ptokens *structs.PTokens) (rerr error) {
	err, client, _ := common.PrepareTokensAndClient(r, ptokens, true)
	if err != nil {
		return err
	}
	userinfo, err := client.Get(cfg.GenOAuth.UserInfoURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := userinfo.Body.Close(); err != nil {
			rerr = err
		}
	}()
	data, _ := ioutil.ReadAll(userinfo.Body)
	log.Infof("Ocs userinfo body: %s", string(data))
	if err = common.MapClaims(data, customClaims); err != nil {
		log.Error(err)
		return err
	}
	ncUser := structs.NextcloudUser{}
	if err = json.Unmarshal(data, &ncUser); err != nil {
		log.Error(err)
		return err
	}
	ncUser.PrepareUserData()
	user.Username = ncUser.Username
	user.Email = ncUser.Email
	return nil
}
