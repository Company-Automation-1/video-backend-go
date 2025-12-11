.PHONY: lint lint-fix test build run

# 代码检查
lint:
	golangci-lint run

# 自动修复可修复的问题
lint-fix:
	golangci-lint run --fix

# 运行测试
test:
	go test -v ./...

# 构建项目
build:
	go build -o app.exe .

# 运行项目
run:
	go run main.go

