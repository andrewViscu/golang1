import React from 'react';
import Login from './Login/Login';
import Register from './Register/Register';
import './Auth.css'
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
  } from "react-router-dom";

function Auth() {
  return (
    <div className="Auth">
        <Router>
        <div className="nav__bar">
            <div className="nav__login">
                <Link to='/login'>Login</Link>
            </div>
            <div className="nav__register">
                <Link to='/register'>Register</Link>
            </div>

            
            <Route path='/login' component={Login} />
            <Route path='/register' component={Register} />


        </div>
        </Router>
        

    </div>
  );
}
{/* <Login /> */}
// <Register /> 
export default Auth;
