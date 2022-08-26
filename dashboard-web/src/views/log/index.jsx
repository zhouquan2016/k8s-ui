import { useEffect, useRef } from "react";
import { Terminal } from "xterm";
import { useSearchParams } from "react-router-dom";
import { FitAddon } from "xterm-addon-fit"

export default function Log() {
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
        const websocket = new WebSocket(`ws://localhost:8080/log/${searchParams.get('namespace')}/${searchParams.get('pod')}`);

        websocket.onopen = () => {
            const fitAddon = new FitAddon();
            terminal.loadAddon(fitAddon);
            terminal.open(xtermRef.current);
            fitAddon.fit();
            return () => {
                terminal.dispose();
            }
        }
        websocket.onmessage = (message) => {
            terminal.write(message.data);
        };
        return () => {
            websocket.close();
        }
    }, [searchParams]);

    return <div ref={xtermRef}/>
}