import App from "../App";
import Cart from "../pages/Cart";
import Home from "../pages/Home";
import Login from "../pages/Login";
import JuiceDetails from "../pages/JuiceDetails";
import OrderHistory from "../pages/OrderHistory";
import AdminLayout from "../pages/AdminLayout";
import AdminProducts from "../pages/AdminProducts";
import Signup from "../pages/Signup";
import AdminOrders from "../pages/AdminOrders";

export const routes = [
    {
        path: '/',
        element: <App />,
        children: [
            {index: true, element:<Home />},
            {path: 'cart', element: <Cart />},
            {path: '/juices/:juiceName', element: <JuiceDetails />},
            {path: '/order-history', element: <OrderHistory />}
        ]
    },
    {
        path: `/login`,
        element: <Login />
    },
    {
        path: 'signup',
        element: <Signup />
    },
    {
        path: '/admin',
        element: <AdminLayout />,
        children: [
            {index: true, element:<AdminProducts />},
            {path: 'products', element:<AdminProducts />},
            {path: 'orders', element:<AdminOrders />}
        ]
    }
]