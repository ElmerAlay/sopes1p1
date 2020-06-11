#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/fs.h>

#define BUFESIZE    150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribir informacion de la memoria ram");
MODULE_AUTHOR("Elmer Edgardo Alay Yupe - 201212945");


struct sysinfo inf;

static int escribir_archivo(struct seq_file * archivo, void *v){
    si_meminfo(&inf);
    seq_printf(archivo, " _____________________________________ \n");
    seq_printf(archivo, "|                                     |\n");
    seq_printf(archivo, "| 201212945                           |\n");
    seq_printf(archivo, "| Elmer Edgardo Alay Yupe             |\n");
    seq_printf(archivo, "|_____________________________________|\n");
    seq_printf(archivo, "\n");
    seq_printf(archivo, " RAM Total: \t%li MB\n", (inf.totalram*4)/(1024));
    seq_printf(archivo, " RAM Libre: \t%li MB\n", (inf.freeram*4)/(1024));
    seq_printf(archivo, " % RAM en uso: \t%li %%\n",((inf.totalram-inf.freeram)*100)/inf.totalram);

    return 0;
}

static int al_abrir(struct inode *inode, struct file *file){
    return single_open(file, escribir_archivo, NULL);
}

static struct file_operations operaciones = 
{
    .open = al_abrir,
    .read = seq_read
};

static int iniciar(void){
    proc_create("memo_201212945", 0, NULL, &operaciones);
    printk(KERN_INFO "Carnet: 201212945\n");
    return 0;
}

static void salir(void){
    remove_proc_entry("memo_201212945", NULL);
    printk(KERN_INFO "Curso: Sistemas Operativos 1\n");
}

module_init(iniciar);
module_exit(salir);