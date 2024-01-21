import React from "react";

function Counter(props) {
    const [count, setCount] = React.useState(props.defaultNum);
    return (
        <div>
            <h1>{count}</h1>
            <button onClick={() => setCount(count + 1)}>Click me</button>
        </div>
    );
}

export default Counter;
