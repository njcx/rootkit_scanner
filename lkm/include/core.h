#ifndef CORE_H
#define CORE_H

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/version.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mutex.h>
#include <linux/workqueue.h>
#include <linux/kallsyms.h>
#include <linux/fs.h>
#include <linux/slab.h>
#include <linux/kallsyms.h>
#include <asm/asm-offsets.h> /* NR_syscalls */

#define MAX_BUFFER_SIZE 4096

struct detection_result {
    char module_alerts[MAX_BUFFER_SIZE];
    char syscall_alerts[MAX_BUFFER_SIZE];
    char process_alerts[MAX_BUFFER_SIZE];
    char interrupt_alerts[MAX_BUFFER_SIZE];
    struct mutex lock;
};

extern struct detection_result global_results;
void init_del_workqueue(void);
void exit_del_workqueue(void);

#endif /* CORE_H */
