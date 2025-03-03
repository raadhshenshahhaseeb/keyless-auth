export interface User {
    id?: number;
    google_id?: string;
    email: string;           
    name: string;
    profile_picture?: string;
    created_at?: Date;
}