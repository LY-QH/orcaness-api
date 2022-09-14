# Orcaness API

## 1. 目录结构
```
├── README.md  // 说明
├── app  // 应用
│   └── domain  // 领域
│       ├── platform  // 平台领域
│       │   ├── entity  // 平台实体
│       │   │   ├── wework_corp.go  // 企业微信主体实体，与 platform_wework_corp 表对应
│       │   │   └── wework_user.go  // 企业微信用户实体，与 platform_wework_user 表对应
│       │   ├── infrastructure  // 基础设施，操作数据库的实现，与 repository 对应
│       │   ├── interface  // 接口层，供外部系统调用
│       │   ├── repository  // 数据仓库，定义数据库方法，由基础设施实现数据库操作
│       │   └── service  // 服务层，供内部系统调用
│       └── user  // 用户领域
│           ├── entity  // 用户实体
│           │   └── base.go  // 实体，与 user_base 表对应
│           ├── infrastructure  // 基础设施，操作数据库的实现，与 repository 对应
│           ├── interface  // 接口层，供外部系统调用
│           │   └── user.go  // 接口实现
│           ├── repository  // 数据仓库，定义数据库方法，由基础设施实现数据库操作
│           └── service  // 服务层，供内部系统调用
│               └── user.go  // 服务实现
├── bin  // 存放生产环境编译的二进制包
├── config  // 配置目录
│   ├── app.dev.yaml  // 开发环境配置文件
│   └── app.yaml  // 默认配置
├── go.mod  // 依赖包集合
├── go.sum  // 依赖包 hash 值
├── main.go  // 主文件
├── router.go  // 路由文件
├── runtime  // 运行时日志目录
├── shell  // 脚本目录
│   ├── build-prod.sh  // 编译生产环境二进制包
│   ├── command.go  // go 命令行实现
│   ├── db-model-generator.sh  // 创建数据表模型
│   └── run-dev.sh  // 运行开发环境
└── util  // 工具目录
    └── package.go  // 获取包名，用于路由解析
```

## 2. 数据字典
表名规范：`领域_实体名称`，如：`wework_user`，表示 `wework` 领域下的 `user` 实体

表字段规范：每张表必须包含 `created_at`、`updated_at` 和 `deleted_at` 字段，且 `deleted_at` 字段需要设置普通索引

## .3 创建模型
```
./shell/db-model-generator.sh tablename
```

## 4. 配置
`app.yaml` 为默认配置

`app.dev.yaml` 为开发环境配置

`app.prod.yaml` 为生产环境配置

*注：运行或编译时会根据环境加载对应的配置文件，并与默认配置合并，环境配置优先于默认配置，合并配置最多支持 3 级，超过 3 级会被替换*

## 5. 运行开发环境
```
./shell/run-dev.sh
```

## 6. 编译生产环境二进制包
```
./shell/build-prod.sh
```