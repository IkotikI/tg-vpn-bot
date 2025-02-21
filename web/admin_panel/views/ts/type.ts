

export type TelegramID = number
export type UserID = number
export type ServerID = number

export interface User {
    ID: UserID;
    TelegramID: TelegramID;
    TelegramName: string;
    CreatedAt: Date;
    UpdatedAt: Date;
}

export interface VPNServer {
    ID: ServerID;
    CountryID: number;
    Name: string;
    Protocol: string;
    Host: string;
    Port: number;
    Username: string;
    Password: string;
    CreatedAt: Date;
    UpdatedAt: Date;
}

export interface Subscriptions {
    UserID: UserID;
    ServerID: ServerID;
    SubscriptionStatus: string;
    SubscriptionExpiredAt: Date;
}

export interface Country {
    CountryID: number;
    Name: string;
    Flag: string;
}