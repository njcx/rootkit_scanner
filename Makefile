SRCDIR := lkm
ORIGIN := $(PWD)
obj-m += rootkit_sc_driver.o
rootkit_sc_driver-objs := core.o util.o proc.o module_list.o syscall_hooks.o  interrupt_hooks.o
HEADERS := $(PWD)/lkm/include
ccflags-y += -I$(HEADERS)


.PHONY: all

all: tidy build module


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

module:
	@echo "构建LKM模块:"
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules