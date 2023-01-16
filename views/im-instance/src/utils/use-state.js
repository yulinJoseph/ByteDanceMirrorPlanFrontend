import { useEffect, useState, useRef } from "react";

function useCallbackState(state) {
    const cbRef = useRef();
    const [data, setData] = useState(state);

    useEffect(() => {
        cbRef.current && cbRef.current(data);
    }, [data]);

    return [
        data,
        function (val, callback) {
            cbRef.current = callback;
            setData(val);
        },
    ];
}

export { useCallbackState };
