import React from 'react';
import './Login.css';

function Login() {
  return (
    <div className="Login">
      <div className="login__form">
        <div className="name">
            <b>Вход</b>
        </div>
      <form>
        <label className="label__form__email">
            Email:<br/>
            <input type="text" name="email" />
        </label>
      </form>
      <form>
        <label className="label__form__email">
            Password:<br/>
            <input type="password" name="password" />
        </label>
        
      </form>
        <form>
          <input type="submit" value="Войти" />
        </form>
      </div>
    </div>
  );
}

export default Login;
