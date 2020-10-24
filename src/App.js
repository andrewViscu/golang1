import React from 'react';
import logo from './logo.svg';
import './App.css';
import Auth from './components/Auth'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  Redirect,
  useParams,
  useRouteMatch
} from "react-router-dom";

function App() {
  return (
    <div className="App">
      <Auth />
    </div>
  );
}

export default App;
