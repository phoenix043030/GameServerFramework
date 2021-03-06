package clientmodel

import (
	"net/http"

	"ssf4g/common/crypto"
	"ssf4g/common/http-const"
	"ssf4g/common/tlog"
	"ssf4g/gamedata/resp-data"
	"ssf4g/server/login-srv/common/err-code"
	"ssf4g/server/login-srv/dao-service/database"
	"ssf4g/server/login-srv/dao-service/memcached"
)

func AccountLogin(w http.ResponseWriter, accntname, accntpass, realip string) *tlog.ErrData {
	accountDB, errData := dbmgr.GetLoginDao().FirstOrInitAccount(accntname)

	if errData != nil {
		errMsg := tlog.Error("account login model (%s) err (first init accnt %v).", accntname, errData.Error())

		respdata.BuildRespFailedRetData(w, httpconst.STATUS_CODE_TYPE_SERVER_ERROR, "account database err")

		return errData.AttachErrMsg(errMsg)
	}

	if accountDB == nil || accountDB.AccntId == 0 {
		tlog.Warn("account login model (%s, %s, %s) warn (account not register).", accntname, accntpass, realip)

		respData := map[string]interface{}{
			"err_code": errcode.LOGIN_ERR_CODE_TYPE_NOT_REGISTER,
		}

		respdata.BuildRespSuccessRetData(w, 0, respData)

		return nil
	}

	accntPass := crypto.EncryptSha1Hash(accntname + accntpass)

	if accntPass != accountDB.PassHash {
		tlog.Warn("account login model (%s, %s) warn (accnt pass illegal).", accntpass, accntPass)

		respData := map[string]interface{}{
			"err_code": errcode.LOGIN_ERR_CODE_TYPE_PASS_ILLEGAL,
		}

		respdata.BuildRespSuccessRetData(w, 0, respData)

		return nil
	}

	accountDB.LastIp = realip

	errData = dbmgr.GetLoginDao().SaveAccount(accountDB)

	if errData != nil {
		errMsg := tlog.Error("account login model (%s, %s) err (save accnt %v).", accntname, realip, errData.Error())

		respdata.BuildRespFailedRetData(w, httpconst.STATUS_CODE_TYPE_SERVER_ERROR, "account database err")

		return errData.AttachErrMsg(errMsg)
	}

	accntToken := crypto.EncryptSha1HashTime(accntPass)

	errData = memcachedmgr.GetAccountMemcached().SetAccntToken(accountDB.AccntId, accntToken)

	if errData != nil {
		errMsg := tlog.Error("account login model (%d, %s) err (set accnt token %v).", accountDB.AccntId, accntToken, errData.Error())

		respdata.BuildRespFailedRetData(w, httpconst.STATUS_CODE_TYPE_SERVER_ERROR, "account memcached err")

		return errData.AttachErrMsg(errMsg)
	}

	respData := map[string]interface{}{
		"accnt_token": accntToken,
	}

	respdata.BuildRespSuccessRetData(w, 0, respData)

	return nil
}
