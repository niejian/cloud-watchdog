// api doc

package api

import "testing"

func TestListNodes(t *testing.T) {
	t.Run("获取节点信息", func(t *testing.T) {
		ListNodes()
	})

}
