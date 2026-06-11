import App from "../App";
import Cart from "../pages/Cart";
import Home from "../pages/Home";
import Login from "../pages/Login";

export const routes = [
    {
        path: '/',
        element: <App />,
        children: [
            {index: true, element:<Home />},
            {path: 'cart', element: <Cart />}
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