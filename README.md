# wrike-confluence-sync
wrike의 프로젝트, 작업 데이터를 confluence 페이지로 작성하는 프로그램이에요

## Architecture
![image](https://user-images.githubusercontent.com/101083786/166687843-9185da01-1523-40ee-8d7f-3ffb966bf2eb.png)

## Directory

```bash
wrike-confluence-sync
│
├── pkg
│   ├──confluence          # confluence API module
│   │  └── ...
│   └──wrike               # wrike API moudle
│      └── ...
└── main.go                # 메인 함수
```
