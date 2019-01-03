package z_log

import "testing"

/**

@author jingsong.zhu
@version 2019/1/3 上午10:13
*/
func TestNewDefault(t *testing.T) {
	lg := GetLogger()
	lg.Info("Testing zap logger.")
}
