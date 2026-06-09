export type LoginCredentials = {
    email: string
    password: string
}

export type LoginRes = {
    username: string,
    email: string,
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

export type MeRes = {
  id: string
  username: string
  email: string
  role: string
  created_at: string
  updated_at: string
}