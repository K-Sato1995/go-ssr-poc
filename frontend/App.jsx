import React from "react";

function App(props) {
    console.log("APP rendered", props);
    const [count, setCount] = React.useState(0);
    return (
        <div>Hello World
        {count}
            <button onClick={() => setCount(count + 1)}>Click me</button>
        </div>
    );
}

export default App;
