# go-zero 模型生成工具介绍

## 1、工具安装
### goctl安装
```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
#验证版本
goctl --version
```

### protoc相关安装

#### 快捷安装
```bash
goctl env check --install --verbose --force
```
#### 手动安装：
- 打开https://github.com/protocolbuffers/protobuf/releases
- 下载对应的版本（我这里是windows电脑），下载protoc-28.0-win64.zip
- 解压，并设置环境变量即可
- 使用protoc --version验证
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --version
```

## 2、快速入门
- 1、创建一个目录gozero-learn做为工作目录
- 2、在当前目录打开命令行工具
- 3、运行goctl api new hello01命令
```bash
goctl api new hello01
```
- 4、生成代码如下：
```bash
SHSZ@DESKTOP-H6EKT5N MINGW64 /d/01_workspace/03_study/01_code/go-zero/hello01
$ tree
.
|-- etc
|   |-- hello01-api.yaml
|-- go.mod
|-- hello01.api
|-- hello01.go
|-- internal
    |-- config
    |   |-- config.go
    |-- handler
    |   |-- hello01handler.go
    |   |-- routes.go
    |-- logic
    |   |-- hello01logic.go
    |-- svc
    |   |-- servicecontext.go
    |-- types
        |-- types.go

7 directories, 11 files
```
- 5、进入hello01，运行go mod tidy下载依赖
- 6、在logic目录下的hello01logic.go中写入如下代码：
```go
func (l *Hello01Logic) Hello01(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	resp = &types.Response{
		Name:    req.Name,
		Message: "hello01",
	}
	return
}
```
- 7、将hello01.go改为main.go并运行main函数
- 8、访问http://localhost:8888/form/you
  - 修改hello01.api
    ```
    type Request {
        Name string `path:"name"`
    }
    ```
  - 重新生成代码，运行命令 goctl api go --api hello01.api --dir .
    - --api：指定api文件
    - --dir：指定go文件生成的目录
  - 重新运行访问即可

## 3、数据库文件生成
- 1、创建model/user目录
- 2、在model/user目录下创建user.sql文件
```mysql
CREATE TABLE `user`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `register_time` datetime NOT NULL,
  `last_login_time` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
```
- 3、在model目录下运行命令 goctl model mysql ddl --src user.sql --dir .
```bash
goctl model mysql ddl --src user.sql --dir .
```

- 4、生成代码如下：
```bash
SHSZ@DESKTOP-H6EKT5N MINGW64 /d/01_workspace/03_study/01_code/go-zero/testserver/model/user
$ tree
.
|-- user.sql
|-- usermodel.go
|-- usermodel_gen.go
`-- vars.go

0 directories, 4 files
```
## 4、goctl-swagger的使用

- 1、安装goctl-swagger
```bash
go install github.com/zeromicro/goctl-swagger@latest
```
- 2、生成swagger文档（.json文件）
```bash
goctl api plugin -plugin goctl-swagger="swagger -filename task.json" -api task.api -dir .
```

