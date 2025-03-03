#ifndef MODULE_LIST_H
#define MODULE_LIST_H
#include "core.h"

void analyze_modules(void);
const char *find_hidden_module(unsigned long addr);

#endif /* MODULE_LIST_H */
