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
		Content:  "docker镜像迁移我们使用的工具是syncer，项目地址：https://github.com/MR5356/syncer ， 并且这个工具支持多对多的镜像仓库迁移\n\n## 安装syncer\n通过下载对应系统的二进制文件进行安装：[点击下载](https://github.com/MR5356/syncer/releases)\n\n也可以通过源码进行安装，前提是有golang运行环境：\n```shell\ngit clone https://github.com/MR5356/syncer.git\ncd syncer\nmake all\n```\n\n## 使用syncer进行镜像迁移\n安装成功后可以使用以下命令获取命令的help信息：\n```shell\n[root@toodo ~] ./syncer -h\n\nUsage:\n  syncer [command]\n\nAvailable Commands:\n  completion  Generate the autocompletion script for the specified shell\n  git         A git repo sync tool\n  help        Help about any command\n  image       A registry image sync tool\n\nFlags:\n  -d, --debug     enable debug mode\n  -h, --help      help for syncer\n  -v, --version   version for syncer\n\nUse \"syncer [command] --help\" for more information about a command.\n\n```\n\n```shell\n[root@toodo ~] ./syncer image -h\n\nA registry image sync tool implement by Go.\n\nComplete code is available at https://github.com/Mr5356/syncer\n\nUsage:\n  syncer image [flags]\n\nFlags:\n  -c, --config string   config file path\n  -d, --debug           enable debug mode\n  -h, --help            help for image\n  -p, --proc int        process num (default 10)\n  -r, --retries int     retries num (default 3)\n  -v, --version         version for image\n```\n\n配置文件支持`yaml`格式和`json`格式，以`yaml`格式为例：\n```yaml\n# 仓库认证信息\nauth:\n  registry.cn-hangzhou.aliyuncs.com:\n    username: your_name\n    password: your_password\n    # http仓库可设置为true\n    insecure: false\n  docker.io:\n    username: your_name\n    password: your_password\n    insecure: false\n# 镜像同步任务列表\nimages:\n  # 该镜像的所有标签将会进行同步\n  registry.cn-hangzhou.aliyuncs.com/toodo/alpine: registry.cn-hangzhou.aliyuncs.com/toodo/test\n  # 该镜像会同步到目标仓库，并使用新的tag\n  alpine@sha256:1fd62556954250bac80d601a196bb7fd480ceba7c10e94dd8fd4c6d1c08783d5: registry.cn-hangzhou.aliyuncs.com/toodo/test:alpine-latest\n  # 该镜像会同步至多个目标仓库，如果目标镜像没有填写tag，将会使用源镜像tag\n  alpine:latest:\n    - hub1.test.com/library/alpine\n    - hub2.test.com/library/alpine\n# 最大并行数量\nproc: 3\n# 最大失败重试次数\nretries: 3\n```\n\n使用配置文件运行镜像迁移工具开始镜像迁移：\n```shell\n[root@toodo ~] ./syncer image -c config.yaml\n```\n",
		Likes:    12312,
		Reads:    1021232,
		Category: []string{"云原生", "容器", "镜像迁移"},
	})
}

func (c *Controller) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/blog")
	api.GET("/:id", c.handleBlogDetail)
}
