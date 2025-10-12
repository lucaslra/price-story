export interface User {
  id: string
  email: string
  password_hash: string
}

export interface Product {
  id: string
  product_name: string
  product_url: string
  product_image_url: string
  product_description: string
  created_by_user: User
  updated_by_user?: User | null
  created_datetime: string
  updated_datetime: string
}
