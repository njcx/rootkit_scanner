SRCDIR := $(PWD)
obj-m += rootkit_sc_driver.o
rootkit_sc_driver-objs := core.o util.o proc.o module_list.o syscall_hooks.o  interrupt_hooks.o
HEADERS := $(PWD)/include
ccflags-y += -I$(HEADERS)

.PHONY: all
all: module cp clean

module:
	@echo "构建LKM模块:"
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

cp:
	cp $(PWD)/rootkit_sc_driver.ko $(PWD)/../rootkit_sc_driver.ko

clean:
	@echo "清理LKM模块:"
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean