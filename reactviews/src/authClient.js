import decodeJwt from 'jwt-decode';
import request from 'superagent';
import { MyConfig } from './MyConfig';

const authClient = (authUrl) => ({
    login: async () => {console.log('login'); return request.get(authUrl);},
    logout: async () => {
        localStorage.removeItem('gotoken');
        return Promise.resolve();
    },
    checkAuth: async () => {console.log('checkAuth');
        const decodedToken = decodeJwt(localStorage.getItem('gotoken'));
        console.log('== getIdentity() : ' + JSON.stringify({ id: decodedToken.id, fullName: decodedToken.id, avatar: '' }) );
        return Promise.resolve();},
    checkError: async (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            localStorage.removeItem('gotoken');
            return Promise.reject();
        }
        return Promise.resolve();
    },
    getPermissions: params => Promise.resolve(),
    getIdentity: async () => {
        console.log('**getIdentity**');
        if ( localStorage.getItem('gotoken') !== null ) {
            const decodedToken = decodeJwt(localStorage.getItem('gotoken'));
            //console.log('== getIdentity() : ' + JSON.stringify({ id: decodedToken.id, fullName: decodedToken.id, avatar: '' }) );
            return { id: decodedToken.id, fullName: decodedToken.id, avatar: '' };
        }
        return { id: '', fullName: '', avatar: ''};
    },
    handleCallback: async () => {
        console.log('**handleCallback**');
        var match = window.location.href.match(/\?(.*)$/);
        console.log(match[1]);
        const token = match[1]; 
        localStorage.setItem('gotoken', token);
        window.location.href = MyConfig.BASE_PATH;
        return Promise.resolve();
    },
});

export default authClient;

