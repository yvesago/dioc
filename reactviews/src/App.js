import React from 'react';
import { fetchUtils, Admin, Resource, Delete } from 'admin-on-rest';
//import './App.css';

import AlertIcon from 'material-ui/svg-icons/action/assessment';
import SurveyIcon from 'material-ui/svg-icons/device/wifi-tethering';
import AgentIcon from 'material-ui/svg-icons/notification/wifi';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import myTheme from './myTheme';

import { Dashboard } from './Dashboard';
import { AlertList, AlertEdit } from './alert';
import { SurveyList, SurveyEdit, SurveyCreate } from './survey';
import { AgentList, AgentEdit, AgentCreate } from './agent';
import { BoardEdit } from './boarddocs';

import { MyConfig } from './MyConfig';
import mySimpleRest from './myJsonRestServer';

import authLoginPage from './authLoginPage';
import authClient from './authClient';
import customRoutes from './customRoutes';

const httpClient = (url, options = {}) => {
    if (!options.headers) {
        options.headers = new Headers({ Accept: 'application/json' });
    }
    // add your own headers here
    //options.headers.set('X-MyToken', MyConfig.API_KEY );
    const token = localStorage.getItem('token');
    options.headers.set('Authorization', `Bearer ${token}`);
    return fetchUtils.fetchJson(url, options);
};


const restClient = mySimpleRest( MyConfig.API_URL  + '/admin/api/v1', httpClient);


const App = () => (
    <Admin title='Distributed IOC manager' theme={getMuiTheme(myTheme)} 
        loginPage={authLoginPage} authClient={authClient(MyConfig.AUTH_URL)}
        customRoutes={customRoutes} dashboard={Dashboard} 
        restClient={restClient}>
        <Resource name="alertes" list={AlertList} edit={AlertEdit} remove={Delete} icon={AlertIcon} />
        <Resource name="surveys" list={SurveyList}  edit={SurveyEdit} create={SurveyCreate} remove={Delete} icon={SurveyIcon} />
        <Resource name="agents" list={AgentList}  edit={AgentEdit} create={AgentCreate} remove={Delete} icon={AgentIcon} />
        <Resource name="board" edit={BoardEdit} />
    </Admin>
);

export default App;

