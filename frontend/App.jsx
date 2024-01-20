import React from "react";
import ReactDOMServer from 'react-dom/server';

function App() {
  const [count, setCount] = React.useState(initialCount);
  return (
    <div className="home">
      <h1>Go + React</h1>
      {count}
    </div>
  );
}

export default App;
