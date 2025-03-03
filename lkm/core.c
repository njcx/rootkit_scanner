//https://github.com/nbulischeck/tyton.git

#include "core.h"
#include "util.h"
#include "proc.h"
#include "logger.h"
#include "module_list.h"
#include "syscall_hooks.h"
#include "interrupt_hooks.h"

static int timeout = 5;
unsigned long *idt = NULL; /* IDT Table */
unsigned long *sct = NULL; /* Syscall Table */
int (*ckt)(unsigned long addr) = NULL; /* Core Kernel Text */

static void work_func(struct work_struct *dummy);
static DECLARE_DELAYED_WORK(work, work_func);
struct detection_result global_results;
static struct proc_dir_entry *proc_entry;

static void clear_results(void) {
    mutex_lock(&global_results.lock);
    memset(global_results.module_alerts, 0, sizeof(global_results.module_alerts));
    memset(global_results.syscall_alerts, 0, sizeof(global_results.syscall_alerts));
    memset(global_results.process_alerts, 0, sizeof(global_results.process_alerts));
    memset(global_results.interrupt_alerts, 0, sizeof(global_results.interrupt_alerts));
    mutex_unlock(&global_results.lock);
}


static void execute_analysis(void) {
    clear_results();
    analyze_modules();
    analyze_syscalls();
    analyze_processes();
    analyze_interrupts();
}


static int results_show(struct seq_file *m, void *v) {
    mutex_lock(&global_results.lock);

    seq_printf(m, "=== Module Analysis Results ===\n\n%s\n", global_results.module_alerts);
    seq_printf(m, "=== Syscall Analysis Results ===\n\n%s\n", global_results.syscall_alerts);
    seq_printf(m, "=== Process Analysis Results ===\n\n%s\n", global_results.process_alerts);
    seq_printf(m, "=== Interrupt Analysis Results ===\n\n%s\n", global_results.interrupt_alerts);

    mutex_unlock(&global_results.lock);
    return 0;
}


static int results_open(struct inode *inode, struct file *file) {
    return single_open(file, results_show, NULL);
}

#if LINUX_VERSION_CODE >= KERNEL_VERSION(5,6,0)
static const struct proc_ops proc_fops = {
    .proc_open = results_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};
#else
static const struct file_operations proc_fops = {
    .owner = THIS_MODULE,
    .open = results_open,
    .read = seq_read,
    .llseek = seq_lseek,
    .release = single_release,
};
#endif


static void work_func(struct work_struct *dummy){
	execute_analysis();
	schedule_delayed_work(&work,
		round_jiffies_relative(timeout*60*HZ));
}

void init_del_workqueue(void){
	schedule_delayed_work(&work, 0);
}

void exit_del_workqueue(void){
	cancel_delayed_work_sync(&work);
}

static int init_kernel_syms(void){
	idt = (void *)lookup_name("idt_table");
	sct = (void *)lookup_name("sys_call_table");
	ckt = (void *)lookup_name("core_kernel_text");

	if (!idt || !sct || !ckt)
		return -1;

	return 0;
}

static int __init init_mod(void) {
    INFO("Inserting Module\n");
    mutex_init(&global_results.lock);
    proc_entry = proc_create("rsc_lkm", 0444, NULL, &proc_fops);
    if (!proc_entry) {
        ERROR("Failed to create proc entry\n");
        return -1;
    }

    if (init_kernel_syms() < 0) {
        ERROR("Failed to lookup symbols\n");
        remove_proc_entry("rsc_lkm", NULL);
        return -1;
    }

    init_del_workqueue();
    return 0;
}


static void __exit exit_mod(void) {
    INFO("Exiting Module\n");
    exit_del_workqueue();
    remove_proc_entry("rsc_lkm", NULL);
}


MODULE_AUTHOR("Nick Bulischeck <nbulisc@clemson.edu>");
MODULE_DESCRIPTION("Linux Kernel-Mode Rootkit Hunter for 4.4.0-31+.");
MODULE_LICENSE("GPL");

module_param(timeout, int, 0);
module_init(init_mod);
module_exit(exit_mod);
