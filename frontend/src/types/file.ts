export interface File {
    id: number;
    extension: string;
    original_name: string;
    format: string;
    size: number;
    size_processed: number;
    status: string;
    status_label: string;
    created_at: string;
    updated_at: string;
}