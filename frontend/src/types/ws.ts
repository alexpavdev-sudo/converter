export interface RegisterAck {
    guest_id: string;
    status: string;
}

export interface NotificationDto {
    id: number;
    detail: string;
    guest_id: number;
    type: 1 | 2;
    created_at: string;
}

export interface WSMessage {
    type: string;
    payload: RegisterAck | NotificationDto | unknown;
}