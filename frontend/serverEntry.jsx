import React from 'react';
import ReactDOMServer from 'react-dom/server';

function App(props) {
    return (
        <div>Hello World</div>
    );
}

export function renderApp(props) {
    return ReactDOMServer.renderToString(<App {...props} />);
}
