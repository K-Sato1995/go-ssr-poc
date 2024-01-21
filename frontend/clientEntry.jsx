import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App'; // Adjust the import path as necessary

const root = ReactDOM.hydrateRoot(
    document.getElementById('app'),
    <App {...(window.APP_PROPS || {})} />
);
