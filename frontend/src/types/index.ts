export type LoginCredentials = {
    email: string
    password: string
}

export type LoginRes = {
    username: string,
    email: string,
    token: string,
    refreshToken: string
}

type UUID = string & { __brand: 'UUID' };

export type SigunpCredentials = {
    email: string,
    username: string,
    password: string
}

export type SignupRes = {
    id: UUID
    createdAt: string
    updatedAt: string
    email: string,
    username: string
}