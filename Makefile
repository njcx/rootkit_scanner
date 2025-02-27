SUBDIRS = lkm

.PHONY: all
all: tidy build $(SUBDIRS)


tidy:
	@echo "整理go依赖包:"
	go mod tidy


build:
	@echo "构建go项目，生成ELF文件:"
	go build -ldflags '-extldflags "-static"'


clean:
	@echo "清理go项目："
	sudo rm -rf rootkit_scanner
	go clean -cache


$(SUBDIRS):
	$(MAKE) -C $@