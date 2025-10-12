export interface User {
  id: string
  email: string
  password_hash: string
  preferred_currency: string
  decimal_places: number
  thousand_separator: string
  currency_symbol_placement: 'before' | 'after'
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
