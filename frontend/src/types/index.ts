export type LoginCredentials = {
    email: string
    password: string
}

export type LoginRes = {
    username: string,
    email: string,
    token: string
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

export type JuiceData = {
    id: string
    name: string
    description: string
    price: number
    stock: number
    image_url: string
    avg_rating: number
    reviews_count: number
}

export type JuiceUpdateParams = {
  name: string,
  description: string,
  price: number,
  stock: number
}

export type OrderItem = {
  name: string
  quantity: number
}

export type OrderDetail = {
  order_id: string
  total: number
  created_at: string
  items: OrderItem[]
}

export type AdminOrder = {
  ID: string
  Total: number
  CreatedAt: string
  UserID: string
  Username: string
}

export type userData = {
  id: string
  username: string
  email: string
  role: string
  created_at: string
  updated_at: string
}