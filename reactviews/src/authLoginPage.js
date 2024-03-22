import { useState } from 'react';
import { MyConfig } from './MyConfig';

const authLoginPage  = ({ theme }) => {

    return (
        <div>
            <h1>DIOC</h1>
            <p>Login with <a href={MyConfig.AUTH_URL}>CAS</a></p>
        </div>
    );
};

export default authLoginPage;
