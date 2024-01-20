import React from 'react';
import { renderToString } from "react-dom/server.browser";

function App(props) {
    return (
        <div>Hello World</div>
    );
}

function renderApp(props) {
    return renderToString(<App {...props} />);
}

// Attach renderApp to the global object
globalThis.renderApp = renderApp;