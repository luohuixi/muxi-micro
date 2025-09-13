# 项目结构
```gas
helloworld/
├── api/
│   ├── local/               存放本地proto
│   │   └── v1/
│   └── third_party/         存放第三方proto
├── cmd/                     主函数所在地
├── configs/                 读取配置
└── internal/
    ├── infrastructure/      基础设施层（提供DB,redisDB,logger等，中间件等）
    ├── repository/          仓储层（配合自动curd使用）
    ├── server/              grpc注册
    ├── service/             服务层（具体业务逻辑）
    └── wire/                wire依赖注入
```

