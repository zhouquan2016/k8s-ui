import './App.css';
import { useRoutes } from 'react-router-dom';
import Index from './views/index/Index';
import Term from './views/term';
import Log from './views/log';
import Editor from './views/editor';

function App() {
  const elements = useRoutes([
    {
      path: '/',
      element: <Index />,
      children: [
        
      ]
    },
    {
      path: '/term',
      element: <Term />,
    },
    {
      path: "/log",
      element: <Log />
    },
    {
      path: '/editor',
      element: <Editor />,
    }
  ])
  return elements;
}

export default App;
