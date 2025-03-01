export interface User {
    id?: number;
    google_id?: string;
    email: string;            // Added required email field
    name: string;
    profile_picture?: string;
    created_at?: Date;
}