package data_splitting

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplitFunction(t *testing.T) {

	// 测试数据
	testString := "没有共产党就没有新中国没有共产党就没有新中国共产党，辛劳为民族共产党他一心救中国他指给了人民解放的道路他领导中国走向光明他坚持了抗战八年多他改善了人民的生活他建设了敌后根据地他实行了民主好处多没有共产党就没有新中国没有共产党就没有新中国（间奏）没有共产党就没有新中国没有共产党就没有新中国共产党，辛劳为民族共产党他一心救中国他指给了人民解放的道路他领导中国走向光明他坚持了抗战八年多他改善了人民的生活他建设了敌后根据地他实行了民主好处多没有共产党就没有新中国没有共产党就没有新中国"
	//testString := "2024-10-15 16:21:23.874 DEBUG [http-nio-8601-exec-2] c.c.u.c.c.OGremlinConnection {.quot;Message.quot;:.quot;sql(g.V().hasLabel(.quot;operatingSystem.quot;).has(.quot;id.quot;, P.within(assetIdList)).project(.quot;id.quot;,.quot;systemType.quot;,.quot;hostName.quot;,.quot;isVirtual.quot;,.quot;ipList.quot;,.quot;empeeAcct.quot;).by(.quot;id.quot;).by(.quot;systemType.quot;).by(.quot;hostName.quot;).by(.quot;isVirtual.quot;).by(.out().hasLabel(.quot;ip.quot;).values(.quot;ip.quot;).fold().map{it.get().join(.quot;,.quot;)}).by(.out().hasLabel(.quot;account.quot;).values(.quot;name.quot;).fold().map{it.get().join(.quot;,.quot;)})) parameters(assetIdList=[9e619f97-fb9b-47af-9433-0e9e0618ff77, 9689600b-8d6a-40c9-a823-ab315ea684bb, d96e2cb4-d52c-4abc-a26e-265acb4b05bd],).quot;,.quot;User.quot;:.quot;root.quot;,.quot;IP.quot;:.quot;134.64.110.143,134.95.172.12:36296.quot;,.quot;EmpeeAcct.quot;:.quot;admin.quot;,.quot;Name.quot;:.quot;cmdb.quot;}"
	chunkSize := int64(11600)

	// 调用split方法
	result, err := split(testString, chunkSize, 1000000)
	if err != nil {
		t.Fatalf("分割字符串失败: %v", err)
	}

	// 打印结果
	fmt.Printf("分割结果（共%d段）:", len(result))
	for _, item := range result {
		fmt.Printf("索引:%d, 大小:%d, 内容:%s，offset：%d\n", item.Index, item.Size, item.Message, item.Offset)
	}

	// 验证分割结果
	if len(result) == 0 {
		t.Error("分割结果为空")
	}

	// 验证每段的大小不超过指定大小
	for i, item := range result {
		if item.Size > chunkSize {
			t.Errorf("第%d段大小超过限制: %d > %d", i+1, item.Size, chunkSize)
		}
		if item.Index != i+1 {
			t.Errorf("索引不正确: 期望%d, 实际%d", i+1, item.Index)
		}
	}

	// 验证所有分块拼接后与原字符串一致
	var reconstructed strings.Builder
	for _, item := range result {
		reconstructed.WriteString(item.Message)
	}
	if reconstructed.String() != testString {
		t.Error("分块拼接后与原字符串不一致")
	}
}
