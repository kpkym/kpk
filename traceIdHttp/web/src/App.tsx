import {useState} from "react";
import * as Diff from "diff";
import service from './axios.js'

import "./App.css";

function App() {
    const [data, setData] = useState({
        logRequestHeader: {},
        logRequestBody: {},
        logResponseBody: {},

        retryResponseHeader: {},
        retryResponseBody: {},
    });
    const [traceId, setTraceId] = useState("");
    const [useDiff, setUseDiff] = useState(true);

    async function click(type: string) {
        let uri = '';
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        let apply = (e: unknown) => {}

        if (type === "log") {
            uri = 'getEsTraceLog';
            apply = (e) => {
                data.logRequestHeader = e.requestHeader || {};
                data.logRequestBody = e.requestBody || {};
                data.logResponseBody = e.responseBody || {};

                data.retryResponseHeader = {};
                data.retryResponseBody = {};
                return {...data};
            };
        }
        if (type === "retry") {
            uri = 'retryTraceReq';
            await click('log');
            apply = (e) => {
                data.retryResponseHeader = e.responseHeader || {};
                data.retryResponseBody = e.responseBody || {};
                return {...data};
            };
        }

        await service.get(`/api/${uri}/${traceId}`)
            .then(e => e.data)
            .then(e => {
                apply(e);
                setData({...data});
            });
    }

    return (
        <>
            <input value={traceId} onChange={(e) => setTraceId(e.target.value)} />
            <button onClick={() => click("log")}>日志</button>
            <button onClick={() => click("retry")}>重试</button>
            <button onClick={() => setUseDiff(!useDiff)}>DIFF {useDiff ? 'ON' : 'OFF'}</button>
            <hr />
            <div className="grid-container">
                <div className="grid-item column-1">
                    <pre>{jsonPretty(data.logRequestHeader)}</pre>
                </div>
                <div className="grid-item column-2">
                    <div className="column-2-up">
                        <pre>{jsonPretty(data.logRequestBody)}</pre>
                    </div>
                    <hr />
                    <div className="column-2-down">
                        <pre>{jsonPretty(data.retryResponseHeader)}</pre>
                    </div>
                </div>
                <div className="grid-item column-3">
                    <pre>{jsonPretty(data.logResponseBody)}</pre>
                </div>
                <div className="grid-item column-4">
                    <pre>{useDiff ? diff(data.logResponseBody, data.retryResponseBody) : jsonPretty(data.retryResponseBody)}</pre>
                </div>
            </div>
        </>
    );
}

function jsonPretty(data: object): string {
    if (Object.keys(data).length === 0) {
        return "";
    }

    return JSON.stringify(data, null, 2);
}

function diff(one: object, other: object) {
    if (Object.keys(other).length === 0) {
        return <pre />;
    }
    const oneStr = jsonPretty(one);
    const otherStr = jsonPretty(other);
    const diffArr = Diff.diffLines(oneStr, otherStr);

    const mappedNodes = diffArr.map((group, index) => {
        const { value, added, removed } = group;
        const color = added ? "#2dba4e" : removed ? 'red' : 'currentcolor';

        return <span key={index} style={{ color }}>{value}</span>;
    });
    return <>{mappedNodes}</>;
}

export default App;
