import { ref, type Ref } from 'vue';
import type { WSMessage } from '@/types/ws';

const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const wsUrl = `${wsProtocol}//${window.location.host}/ws`;

// Глобальное состояние
const isConnected = ref(false);
let ws: WebSocket | null = null;
const listeners: Array<(msg: WSMessage) => void> = [];

let connectionInProgress = false;

export function useWebSocket() {
    const connect = () => {
        if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
            return;
        }

        connectionInProgress = true;

        ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            isConnected.value = true;
            connectionInProgress = false;
        };

        ws.onmessage = (event: MessageEvent) => {
            try {
                const message: WSMessage = JSON.parse(event.data);
                listeners.forEach(cb => cb(message));
            } catch (error) {
            }
        };

        ws.onerror = (error: Event) => {
            connectionInProgress = false;
        };

        ws.onclose = (event: CloseEvent) => {
            isConnected.value = false;
            connectionInProgress = false;
        };
    };

    const disconnect = () => {
        if (ws) {
            ws.close();
            ws = null;
            isConnected.value = false;
        }
        listeners.length = 0; // удаляем всех слушателей
    };

    /**
     * Регистрирует обработчик входящих сообщений.
     * Возвращает функцию для отписки.
     */
    const onMessage = (callback: (msg: WSMessage) => void): (() => void) => {
        listeners.push(callback);
        return () => {
            const idx = listeners.indexOf(callback);
            if (idx > -1) listeners.splice(idx, 1);
        };
    };

    return {
        isConnected,
        connect,
        disconnect,
        onMessage,
    };
}