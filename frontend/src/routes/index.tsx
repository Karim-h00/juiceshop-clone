import App from "../App";
import Cart from "../pages/Cart";
import Home from "../pages/Home";
import Login from "../pages/Login";
import JuiceDetails from "../pages/JuiceDetails";
import OrderHistory from "../pages/OrderHistory";

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
        element: <Login />
    }
]