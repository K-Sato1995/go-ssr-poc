import React from 'react';
import { renderToString } from "react-dom/server.browser";

function App(props) {
    const [count, setCount] = React.useState(0);
    return (
        <div>Hello World
        {count}
            <button onClick={() => setCount(count + 1)}>Click me</button>
        </div>
    );
}

function renderApp(props) {
    return renderToString(<App {...props} />);
}

globalThis.renderApp = renderApp;
