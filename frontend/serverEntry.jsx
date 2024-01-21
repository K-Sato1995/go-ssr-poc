import React from 'react';
import { renderToString } from "react-dom/server.browser";
import App from './App'; // Adjust the import path as necessary

function renderApp(props) {
    return renderToString(<App {...props} />);
}

globalThis.renderApp = renderApp;
