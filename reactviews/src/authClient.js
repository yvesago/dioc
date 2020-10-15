import decodeJwt from 'jwt-decode';
import request from 'superagent';

const authClient = (authUrl) => ({
    login: () => {return request.get(authUrl);},
    logout: () => {
        localStorage.removeItem('token');
        return Promise.resolve();
    },
    checkAuth: () => Promise.resolve(),
    checkError: (error) => {
        const status = error.status;
        if (status === 401 || status === 403) {
            localStorage.removeItem('token');
            return Promise.reject();
        }
        return Promise.resolve();
    },
    getPermissions: params => Promise.resolve(),
    getIdentity: () => {
        if ( localStorage.getItem('token') !== null ) {
            const decodedToken = decodeJwt(localStorage.getItem('token'));
            //console.log('== getIdentity() : ' + JSON.stringify({ id: decodedToken.id, fullName: decodedToken.id, avatar: '' }) );
            return { id: decodedToken.id, fullName: decodedToken.id, avatar: '' };
        }
        return { id: '', fullName: '', avatar: ''};
    },
});

export default authClient;

