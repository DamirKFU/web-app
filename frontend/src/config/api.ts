export const API_CONFIG = {
  baseURL: 'http://localhost:8080',
  endpoints: {
    csrf: '/api/get-csrf-token/',
    login: '/api/v1/auth/login/',
    register: '/api/v1/auth/register/',
    products: '/api/catalog/products/',
    categories: '/api/catalog/categories/',
    colors: '/api/catalog/colors/',
    productDetail: (id: number) => `/api/catalog/products/${id}/`,
    cart: {
      add: '/api/catalog/cart/add/',
      get: '/api/catalog/cart/',
      update: (id: number) => `/api/catalog/cart/item/${id}/`,
    },
    getCSRFToken: '/api/get-csrf-token/',
    orders: '/api/catalog/orders/',
    orderDetail: (id: number) => `/api/catalog/orders/${id}/`,
    support: {
      chats: '/api/support/chats/',
      chatDetail: (id: number) => `/api/support/chats/${id}/`,
      chatMessages: (id: number) => `/api/support/chats/${id}/messages/`,
      invite: (id: number) => `/api/staff/support/chats/${id}/invite/`,
    },
    users: {
      search: '/api/users/search/',
    },
    staff: {  
      chats: '/api/staff/support/chats/',
    },
  },
} as const;
