import React, { Component } from 'react';
import { connect } from 'react-redux';
import { userLogin } from 'react-admin';
import { MyConfig } from './MyConfig';

class authLoginPage extends Component {
    /*submit = (e) => {
        e.preventDefault();
        // gather your data/credentials here
        const credentials = { };

        // Dispatch the userLogin action (injected by connect)
        this.props.userLogin(credentials);
    }*/

    render() {
        return (
            <div>
                <h1>DIOC</h1>
                <p>Login with <a href={MyConfig.AUTH_URL}>CAS</a></p>
            </div>
        );
    }
}

export default connect(undefined, { userLogin })(authLoginPage);

