import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/mode-xml";
import "ace-builds/src-noconflict/mode-properties";
import "ace-builds/src-noconflict/mode-text";
import "ace-builds/src-noconflict/mode-xml";
import "ace-builds/src-noconflict/mode-html";
import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-cloud9_day";
import "ace-builds/src-noconflict/theme-cloud9_night";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/theme-github";

import "ace-builds/src-noconflict/ext-language_tools";
import { useEffect, useRef, useState} from "react";
import { useSearchParams, useLocation, useNavigate } from "react-router-dom"

const baseUrl = "http://localhost:8080";

function EditorKey({ configmap, handlerUpdate, handlerRemoveKey, toggleEditable }) {
    const { key, data, editable, isFold } = configmap;
    const [newValue, setNewValue] = useState(data)
    const [fold, setFold] = useState(isFold)
    const theme = "cloud9_night"
    let mode = key.lastIndexOf(".") !== -1 ? key.substring(key.lastIndexOf(".") + 1).toLocaleLowerCase() : "text"
    if (mode === "yml") {
        mode = "yaml"
    }
    function onChange(newValue) {
        setNewValue(newValue)
    }
    function handlerEdit() {
        if(editable) {
            //保存
            const obj = {}
            obj[key] = newValue
            handlerUpdate(obj)
        }else {
            //取消编辑
            if (toggleEditable(key)) {
                setFold(false)
            }
        }
    }
    function handlerCancelEdit() {
        if(data !== newValue && !window.confirm("Are you sure you want to delete this change?")) {
            return
        }
        setNewValue(data)
        toggleEditable(key)
    }
    return <div>
        <span id={key}>{key}</span>
        <button onClick={handlerEdit}>{editable ? '保存' : '编辑'}</button>
        <button onClick={handlerCancelEdit} disabled={!editable} >取消</button>
        <button onClick={() => handlerRemoveKey(key)} disabled={editable}>删除</button>
        <button onClick={() => setFold(!fold)}>{fold ? '+' : '-'}</button>
        {
            !fold && <AceEditor
                mode={mode}
                theme={theme}
                width="90%"
                name={key}
                onChange={onChange}
                fontSize={14}
                showPrintMargin={true}
                showGutter={true}
                highlightActiveLine={true}
                value={newValue}
                readOnly={!editable}
                focus={editable}
                isInShadow={!editable}
                setOptions={{
                    useWorker: false,
                    enableBasicAutocompletion: true,
                    enableLiveAutocompletion: true,
                    enableSnippets: false,
                    showLineNumbers: true,
                    tabSize: 2,

                }} />
        }
        

    </div>
}

export default function Editor() {
    const [searchParams] = useSearchParams()
    const [list, setList] = useState([]);
    const keyRef = useRef();
    const localtion = useLocation();
    const navigate = useNavigate()
    const namespace = searchParams.get('namespace')
    const name = searchParams.get('name')
    useEffect(() => {
        fetch(`${baseUrl}/configmap/${namespace}/${name}`).then((resp) => {
            resp.json().then((apiResult) => transferApiResult(apiResult))
        });
    }
        , [namespace, name]);
    function transferApiResult(apiResult) {
        if (apiResult.code === 3000) {
            alert(apiResult.error);
            return
        } else {
            setList(apiResult.data)
        }
        if (!apiResult.data) {
            setList([]);
            return
        }
        apiResult.data.forEach(data => {
            data.isFold = false
        });
        setList(apiResult.data)
    }
    function handlerUpdate(data) {
        fetch(`${baseUrl}/configmap/${namespace}/${name}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        }).then((resp) => {
            resp.json().then((apiResult) => transferApiResult(apiResult))
        }).catch(err => {
            if ("error" in err) {
                alert(err.error)
            } else {
                alert(err)
            }
        });
    }
    function handlerAddKey() {
        const newList = list.slice();
        const key = keyRef.current.value
        console.log(localtion);
        if (!key) {
            alert("请输入key")
            return
        }
        const exists = newList.some((item) => item.key === key)
        if (exists) {
            alert("新增项已存在!");
            navigate(`${localtion.pathname}${localtion.search}#${key}`)
            return
        }
        const unsave = newList.filter((item) => item.editable)[0]
        if (unsave) {
            alert(`请先保存${unsave.key}!`);
            return
        }
        newList.push({
            key,
            data: "",
            editable: true,
            isFold: false,
        });
        setList(newList);
        navigate(`${localtion.pathname}${localtion.search}#${key}`)
        
    }
    function handlerRemoveKey(key) {
        if(window.confirm("Are you sure you want to remove this key?")) {
            fetch(`${baseUrl}/configmap/${namespace}/${name}/${key}`, {
                method: "DELETE",
            }).then((resp) => {
                resp.json().then((apiResult) => transferApiResult(apiResult))
            })
        }
        
    }
    function toggleEditable(key) {
        const unsave = list.filter((item) => item.editable)[0]
        if (unsave && unsave.key !== key) {
            alert(`请先保存${unsave.key}!`);
            return false
        }
        const newList = list.slice().map((item) => {
            if(item.key === key) {
                item.editable = !item.editable
            }
            return item
        })
        setList(newList);
        return true
    }
    return (
        <div>
            <div>
                <input type="text" ref={keyRef} />
                <button onClick={handlerAddKey}>新增一项</button>
            </div>
            <div>
                {
                    list.map(config => <EditorKey configmap={config} key={config.key} 
                        handlerUpdate={handlerUpdate} handlerRemoveKey={handlerRemoveKey}
                        toggleEditable={toggleEditable}
                    />)
                }
            </div>
        </div>

    )
}