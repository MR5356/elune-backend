package blog

import (
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) handleBlogDetail(ctx *gin.Context) {
	response.Success(ctx, &Blog{
		ID:       1,
		Title:    "Docker镜像仓库迁移",
		Desc:     "docker镜像迁移我们使用的工具是syncer",
		Author:   "Elune",
		Content:  "# 一级标题\n## 二级标题\n### 三级标题\n#### 四级标题\n##### 五级标题\n###### 六级标题\n\n## 段落\n这里展示的是段落的样式\n\n这是第二个段落\n\n## 代码块\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Print(\"Elune\")\n}\n```\n\n## 粗体\n粗体展示**我粗了**\n\n## 斜体\n斜体展示*我歪了*\n\n粗体加斜体展示***我又歪又粗***\n\n## 引用\n> 引用了伟大物理学家R.Ma的一句话\n\n> 引用一\n>> 嵌套了引用二\n\n## 有序列表\n1. 嘟嘟嘟\n2. 嘟嘟嘟\n3. 嘟嘟嘟\n\n## 无序列表\n* 嘟嘟嘟\n* 嘟嘟嘟\n* 嘟嘟嘟\n\n## 分割线\n\n-----\n\n## 超链接\n这是一个链接 [Elune](https://docker.ac.cn \"Elune\")\n\n## 图片\n![这是图片](https://docker.ac.cn/logo.svg \"Elune\")",
		Likes:    12312,
		Reads:    1021232,
		Category: []string{"云原生", "容器", "镜像迁移"},
	})
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/blog")
	api.GET("/:id", c.handleBlogDetail)
}
