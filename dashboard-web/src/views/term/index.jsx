import { useEffect } from "react";
import { Terminal } from 'xterm';
import "xterm/css/xterm.css"
import "./index.css"
import { FitAddon } from "xterm-addon-fit"
import { useRef } from "react";
import {useSearchParams} from "react-router-dom";
export default function Term() {
    const xtermRef = useRef();
    const [searchParams] = useSearchParams();
    
    useEffect(() => {
        console.log("useEffect");
        const terminal = new Terminal({
            // 渲染类型
            rendererType: 'dom',
            //   是否禁用输入
            disableStdin: false,
            cursorStyle: 'underline',
            //   启用时光标将设置为下一行的开头
            convertEol: true,
            // 终端中的回滚量
            scrollback: 50,
            fontSize: 18,
            rows: 50,
            // 光标闪烁
            cursorBlink: true,
        });
        const websocket = new WebSocket(`ws://localhost:8080/ws?namespace=${searchParams.get('namespace')}&pod=${searchParams.get('pod')}`);
        function sendData(data) {
            websocket.send(JSON.stringify(data));
        }
        websocket.onopen = () => {
            const fitAddon = new FitAddon();
            terminal.loadAddon(fitAddon);
            terminal.onData((data) => sendData({
                type: "input",
                message: data
            }))
            terminal.onResize(({ cols, rows }) => sendData({
                type: "resize",
                cols,
                rows
            }));
            terminal.open(xtermRef.current);
            fitAddon.fit();
            return () => {
                terminal.dispose();
            }
        } 
        websocket.onmessage = (message) => {
            console.log("message=", message);
            terminal.write(message.data);
        };
        return () => {
            websocket.close();
        }
    }, [searchParams]);
    return (<div ref={xtermRef} className="container"/>);
}