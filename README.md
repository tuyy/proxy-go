# Proxy-go
 표준 라이브러리, 설정파일 기반 프록시 서버

## 설정 예시
```
{
    "default_dest_url": "http://localhost:8081",
    "mappings": [
        {"path": "/test1", "dest_url": "http://127.0.0.1:8081"},
        {"path": "/test2", "dest_url": "http://127.0.0.1:8081"},
        {"path": "/test3", "dest_url": "http://127.0.0.1:8081"}
    ],
    "loggers": [
        {"level": "1", "description": "debug", "filepath": "./logs/debug.log", "format": "[%D] %m", "class": "AsyncLogger"} ,
        {"level": "2", "description": "info", "filepath": "./logs/info.log", "format": "[%D] %m", "class": "AsyncLogger"} ,
        {"level": "5", "description": "access", "filepath": "./logs/access.log", "format": "[%D] %m", "class": "AsyncLogger"} ,
        {"level": "6", "description": "debug", "filepath": "./logs/result.log", "format": "[%D] %m", "class": "AsyncLogger"}
    ]
}
```

* default_dest_url: mappings에 등록되지 않은 PATH로 호출되는 경우 호출하는 URl
* mappings: path로 호출되면 dest_url로 호출