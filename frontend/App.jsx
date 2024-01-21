import React from "react";

function App(props) {
    console.log("APP rendered", props);
    const [count, setCount] = React.useState(0);
    return (
        <div>Hello World
            <h1>タイトル: {props.Name}</h1>
            {count}
            <button onClick={() => setCount(count + 1)}>Click me</button>
        </div>
    );
}

export default App;
