#include "core.h"

#define MESSAGE(level, tag, ...) printk(level tag __VA_ARGS__);
#define INFO(...) MESSAGE(KERN_INFO,    "INFO: ", __VA_ARGS__)
#define ALERT(...) MESSAGE(KERN_ALERT,   "ALERT: ", __VA_ARGS__)
#define WARNING(...) MESSAGE(KERN_WARNING, "WARNING: ", __VA_ARGS__)
#define ERROR(...) MESSAGE(KERN_ERR,     "ERROR: ", __VA_ARGS__)
