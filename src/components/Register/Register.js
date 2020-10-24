import React from 'react';
import './Register.css';

function Register() {
  return (
    <div className="Register">
      <div className="register__form">
        <div className="name">
            <b>Регистрация</b>
        </div>
      <form>
        <label className="label__form">
            Имя:<br/>
            <input type="text" name="name" />
        </label>
      </form>
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
          <input type="submit" value="Регистрация" />
        </form>
      </div>
    </div>
  );
}

export default Register;
