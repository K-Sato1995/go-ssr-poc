import React from "react";
import Counter from "./components/Counter";

function App(props) {
    console.log("APP rendered", props);
    const [count, setCount] = React.useState(0);
    return (
        <div>
            <h1>タイトル: {props.Name}</h1>
            <Counter />
        </div>
    );
}

export default App;
