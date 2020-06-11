#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/sched/signal.h>
#include <linux/tty.h>
#include <linux/version.h>
#include <linux/kthread.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/fs.h>

#define MAX_BUF_SIZE 1024

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribir informacion de los procesos del cpu");
MODULE_AUTHOR("Elmer Edgardo Alay Yupe - 201212945");

struct list_head *list;
struct task_struct *task_child;
struct task_struct *iter;

static int escribir_archivo(struct seq_file * archivo, void *v)
{
    seq_printf(archivo, " _____________________________________ \n");
    seq_printf(archivo, "|                                     |\n");
    seq_printf(archivo, "| 201212945                           |\n");
    seq_printf(archivo, "| Elmer Edgardo Alay Yupe             |\n");
    seq_printf(archivo, "|_____________________________________|\n");
    seq_printf(archivo, "\n");

    int cnt = 0;
    for_each_process(iter)
    {
        seq_printf(archivo, "Padre PID:%d - NOMBRE:%s - Estado:%ld\n", iter->pid, iter->comm, iter->state);

        list_for_each(list, &iter->children)
        {
            task_child = list_entry(list, struct task_struct, sibling);
            seq_printf(archivo, " -Hijo de %s[%d] PID:%d - NOMBRE:%s - Estado:%ld\n", iter->comm, iter->pid, task_child->pid, task_child->comm, task_child->state);
        }

        cnt++;
    }

    seq_printf(archivo, "Número de procesos: %d\n", cnt);
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

static int iniciar(void)
{
    proc_create("cpu_201212945", 0, NULL, &operaciones);
    printk(KERN_INFO "Nombre: Elmer Edgardo Alay Yupe\n");
    return 0;
}

static void salir(void)
{
    remove_proc_entry("cpu_201212945", NULL);
    printk(KERN_INFO "Curso: Sistemas Operativos 1\n");
}

module_init(iniciar);
module_exit(salir);