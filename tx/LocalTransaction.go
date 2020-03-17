package tx

import (
	"cta/util"
	"cta/logs"
	"database/sql"
	"time"
)

type LocalTx struct {
	tx         *sql.Tx
	xid        string
	resourceId string
	branchId   int64

	RetryTimes    int
	RetryInterval time.Duration
}

func (ltx *LocalTx) Commit() bool {
	ok := ltx.lockGlobalResource() // get global lock before committing branch
	if !ok {
		logs.Infof("XID: %s, resourceId: %s, branchId: %d, get global lock failed", ltx.xid, ltx.resourceId, ltx.branchId)
		return false
	}
	success := util.Retry(ltx.RetryTimes, ltx.RetryInterval, func() bool {
		err := ltx.tx.Commit()
		return err == nil
	})
	_ = ltx.unlockGlobalResource()

	if success {
		_ = ltx.reportBranch(true)
	} else {
		logs.Warnf("XID: %s, branchId: %d, commit branch failed", ltx.xid, ltx.branchId)
		_ = ltx.RollBack()
	}
	return success

}

func (ltx *LocalTx) RollBack() bool {
	success := util.Retry(ltx.RetryTimes, ltx.RetryInterval, func() bool {
		err := ltx.tx.Rollback()
		return err == nil
	})
	if !success {
		logs.Warnf("XID: %s, branchId: %d, rollback branch failed", ltx.xid, ltx.branchId)
	}

	ltx.reportBranch(false)
	return success
}

func (ltx *LocalTx) reportBranch(success bool) bool {
	// todo
	ok := util.Retry(ltx.RetryTimes, ltx.RetryInterval, func() bool {
		return false
	})

	if !ok {
		logs.Warnf("XID: %s, branchId: %d, report branch failed", ltx.xid, ltx.branchId)
	}
	return ok
}

func (ltx *LocalTx) lockGlobalResource() bool {
	// todo
	return false
}

func (ltx *LocalTx) unlockGlobalResource() bool {
	// todo
	return false
}
