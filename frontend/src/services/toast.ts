import { ref } from 'vue';

export interface Toast {
    id: number;
    title: string;
    detail: string;
    type: 'success' | 'error';
}

const toasts = ref<Toast[]>([]);
let nextId = 1;

export function toast() {
    const showToast = (
        title: string,
        detail: string,
        type: 'success' | 'error' = 'error'
    ) => {
        const id = nextId++;
        const toast: Toast = {
            id,
            title,
            detail,
            type,
        };

        toasts.value.push(toast);

        // Автоматическое удаление через 10 секунд
        setTimeout(() => {
            removeToast(id);
        }, 10000);
    };

    const removeToast = (id: number) => {
        const index = toasts.value.findIndex((n) => n.id === id);
        if (index !== -1) {
            toasts.value.splice(index, 1);
        }
    };

    return {
        toasts: toasts,
        showToast: showToast,
        removeToast: removeToast,
    };
}