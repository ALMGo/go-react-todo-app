import {useEffect, useState} from 'react'
import logo from './logo.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0);
  const [username, setUsername] = useState("")
  const [todo, setTodo] = useState({})

  useEffect(() => {
    (async () => {
      const formData = new FormData();
      formData.append("username", "user1");
      formData.append("password", "secret");
      const login = await fetch("/api/user/login", {
        method: "POST",
        body: formData,
      })
      const loginRes = await login.json();
      if(loginRes.status === 200) {
        setUsername(username)
      }
      console.log(loginRes);
      const res = await fetch("/api/todo/1");
      const data = await res.json()
      setTodo(data)
      console.log(data)
    })();
  }, []);


  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        {username && <p>Hello: {username}</p>}
        <p>Hello Vite + React!</p>
        <p>
          <button type="button" onClick={() => setCount((count) => count + 1)}>
            count is: {count}
          </button>
        </p>
        <p>
          Edit <code>App.jsx</code> and save to test HMR updates.
        </p>
        <p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
          {' | '}
          <a
            className="App-link"
            href="https://vitejs.dev/guide/features.html"
            target="_blank"
            rel="noopener noreferrer"
          >
            Vite Docs
          </a>
        </p>
        <div>{JSON.stringify(todo)}</div>
      </header>
    </div>
  )
}

export default App
