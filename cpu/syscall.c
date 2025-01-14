#include <elf.h>
#include <sys/auxv.h>

int auxval(){
    unsigned long v = getauxval(AT_CLKTCK);
    return (int)(v);
}