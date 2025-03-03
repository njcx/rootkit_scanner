#ifndef _LOGGER_H
#define _LOGGER_H

#include "core.h"
#define KERN_LOG KERN_INFO "[RKHUNTER] "
#define KERN_ALERT_LOG KERN_ALERT "[RKHUNTER-ALERT] "
#define KERN_WARN_LOG KERN_WARNING "[RKHUNTER-WARNING] "

#define INFO(fmt, args...) \
    do { \
        printk(KERN_LOG fmt, ##args); \
    } while (0)

#define ERROR(fmt, args...) \
    do { \
        printk(KERN_ALERT_LOG fmt, ##args); \
    } while (0)

#define WARNING(fmt, args...) \
    do { \
        printk(KERN_WARN_LOG fmt, ##args); \
    } while (0)


#define ALERT(fmt, args...) \
    do { \
        char tmp_buf[512]; \
        int len; \
        printk(KERN_ALERT_LOG fmt, ##args); \
        len = snprintf(tmp_buf, sizeof(tmp_buf), "[Warning] " fmt, ##args); \
        if (len > 0) { \
            mutex_lock(&global_results.lock); \
            strlcat(global_results.module_alerts, tmp_buf, \
                   sizeof(global_results.module_alerts)); \
            mutex_unlock(&global_results.lock); \
        } \
    } while (0)


#endif /* _LOGGER_H */