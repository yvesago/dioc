import { AUTH_LOGIN, AUTH_LOGOUT, AUTH_ERROR, AUTH_CHECK } from 'admin-on-rest';
import request from 'superagent';

export default (authUrl) => (type, params) => {
    if (type === AUTH_LOGIN) {
        return request.get(authUrl);
    }

    if (type === AUTH_ERROR) {
        const {status} = params;
        
        if (status === 401 || status === 403) {
            localStorage.removeItem('token');
            window.location.href = '/#/login';
            return Promise.reject();
        }

        return Promise.resolve();
    }

    if (type === AUTH_CHECK) {
        return localStorage.getItem('token') ? Promise.resolve() : Promise.reject();
    }

    if (type === AUTH_LOGOUT) {
        localStorage.removeItem('token');
        return Promise.resolve();
    }

    return Promise.resolve();
};