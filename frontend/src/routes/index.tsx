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
import OrderDetails from "../components/OrderDetails";
import AdminUsers from "../pages/AdminUsers";
import Profile from "../pages/Profile";
import AdminAuditLogs from "../pages/AdminAuditLogs";


export const routes = [
    {
        path: '/',
        element: <App />,
        children: [
            {index: true, element:<Home />},
            {path: 'cart', element: <Cart />},
            {path: '/juices/:juiceName', element: <JuiceDetails />},
            {path: '/order-history', element: <OrderHistory />},
            {path: 'profile', element: <Profile />}
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
            {path: 'orders', element:<AdminOrders />},
            {path: 'orders/:id', element:<OrderDetails />},
            {path: 'users', element: <AdminUsers />},
            {path: 'audits', element: <AdminAuditLogs />}
        ]
    }
]